/*CORE*/
import {Inject, Injectable} from '@angular/core';
import {BehaviorSubject, forkJoin, Observable, of, throwError} from 'rxjs';
import {concatMap, map, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import {WEB3} from './web3';
import {Transaction as Web3Tx, Tx} from 'web3/eth/types';
import {Account, TxSignature} from 'web3/eth/accounts';
import {TransactionReceipt} from 'web3/types';
/*SERVICES*/
import {ToastrService} from '../toastr/toastr.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {ABIDefinition} from 'web3/eth/abi';
import {Transaction} from '../../models/transaction.model';
/*UTILS*/
import {objIsEmpty} from '../../utils/functions';
import {ContractAbi} from '../../utils/types';

@Injectable()
export class WalletService {

  private _abi$: BehaviorSubject<ContractAbi> = new BehaviorSubject<ContractAbi>(null);
  private _abi: ContractAbi;

  get abi$() {
    if (!this._abi) {
      return this.getAbi();
    }
    return this._abi$;
  }

  get w3(): Web3 {
    return this._web3;
  }

  constructor(@Inject(WEB3) public _web3: Web3,
              private _toastrService: ToastrService,
              private _commonService: CommonService) {
    if (!this._web3) {
      return;
    }
    const provider = new this._web3.providers.HttpProvider(this._commonService.rpcProvider);
    this._web3.setProvider(provider);
  }

  getAbi(): Observable<ContractAbi> {
    return this._commonService.getAbi().pipe(
      tap((abi: ContractAbi) => {
        this._abi = abi;
        this._abi$.next(abi);
      })
    );
  }

  createAccount(): Account {
    return !!this._web3 ? this._web3.eth.accounts.create() : null;
  }

  isAddress(address: string) {
    return this._web3.utils.isAddress(address);
  }

  sendTx(privateKey: string, tx: Tx): any {
    let from;
    try {
      from = this._web3.eth.accounts.privateKeyToAccount(privateKey);
    } catch (e) {
      return throwError(e);
    }

    const p = this._web3.eth.getTransactionCount(from.address);
    return fromPromise(p).pipe(
      concatMap(nonce => {
        tx.nonce = nonce;
        const p2: Promise<TxSignature> = this._web3.eth.accounts.signTransaction(tx, privateKey);
        return fromPromise(p2);
      }),
      concatMap((signed: TxSignature) => {
        return this.sendSignedTx(signed);
      })
    );
  }

  sendSignedTx(signed: TxSignature): Observable<TransactionReceipt> {
    return fromPromise(this._web3.eth.sendSignedTransaction(signed.rawTransaction));
  }

  getBalance(address: string): Observable<string> {
    try {
      const p = this._web3.eth.getBalance(address);
      return fromPromise(p).pipe(
        map((balance: string) => this._web3.utils.fromWei(balance, 'ether')),
      );
    } catch (e) {
      return throwError(e);
    }
  }

  /**
   * call function
   * @param addr
   * @param abi
   * @param params
   */
  call(addr: string, abi: ABIDefinition, params: any[]): Promise<object> | null {
    try {
      const encoded: string = this._web3.eth.abi.encodeFunctionCall(abi, params);
      return this._web3.eth.call({
        to: addr,
        data: encoded,
      }).then((res: string) => {
        if (!res) {
          throw new Error('Result is empty');
        }
        const decoded: object = this._web3.eth.abi.decodeLog(abi.outputs, res, []);
        if (objIsEmpty(decoded)) {
          throw new Error('Result is empty');
        }
        return decoded;
      });
    } catch (err) {
      throw err;
    }
  }

  /**
   * getting tx from node in case of server haven't processed yet
   * @param txHash
   */
  getTxData(txHash: string): Observable<Transaction> {
    if (!this._web3) {
      return of(null);
    }
    return forkJoin([
      fromPromise<Web3Tx>(this._web3.eth.getTransaction(txHash)),
      fromPromise<TransactionReceipt>(this._web3.eth.getTransactionReceipt(txHash)),
    ]).pipe(
      map((res: [Web3Tx, TransactionReceipt]) => {
        if (!res[0]) {
          return null;
        }
        const tx: Web3Tx = res[0];
        const txReceipt = res[1];
        const finalTx: Transaction = new Transaction();
        finalTx.tx_hash = tx.hash;
        finalTx.value = tx.value;
        finalTx.gas_price = tx.gasPrice;
        finalTx.gas_limit = '' + tx.gas;
        finalTx.nonce = tx.nonce;
        finalTx.input_data = tx.input.replace(/^0x/, '');
        finalTx.from = tx.from;
        finalTx.to = tx.to;
        if (txReceipt) {
          finalTx.block_number = tx.blockNumber;
          finalTx.gas_fee = '' + (+tx.gasPrice * txReceipt.gasUsed);
          finalTx.contract_address = txReceipt.contractAddress;
          finalTx.status = txReceipt.status;
          finalTx.created_at = new Date();
        }
        return finalTx;
      }),
    );
  }

  estimateGas(tx: Tx): Observable<number> {
    return fromPromise(this._web3.eth.estimateGas(tx));
  }
}

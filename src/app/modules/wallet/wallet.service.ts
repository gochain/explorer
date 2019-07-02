/*CORE*/
import {Inject, Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {BehaviorSubject, forkJoin, Observable, of, throwError} from 'rxjs';
import {concatMap, map, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import {WEB3} from './web3';
import {Transaction as Web3Tx, Tx} from 'web3/eth/types';
import {Account, TxSignature} from 'web3/eth/accounts';
import {TransactionReceipt} from 'web3/types';
import Web3Contract from 'web3/eth/contract';
import {ABIDefinition} from 'web3/eth/abi';
/*SERVICES*/
import {ToastrService} from '../toastr/toastr.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Transaction} from '../../models/transaction.model';
/*UTILS*/
import {getDecodedData, objIsEmpty} from '../../utils/functions';
import {ContractAbi} from '../../utils/types';

@Injectable()
export class WalletService {

  private _abi$: BehaviorSubject<ContractAbi> = new BehaviorSubject<ContractAbi>(null);
  private _abi: ContractAbi;

  isProcessing = false;
  isProcessing$: BehaviorSubject<boolean> = new BehaviorSubject(false);

  // ACCOUNT INFO
  account: Account;
  accountBalance$: BehaviorSubject<any> = new BehaviorSubject<any>(null);
  private accountBalance: string;

  receipt: TransactionReceipt;
  contract: Web3Contract;
  selectedFunction: ABIDefinition;
  functionResult: any[][];

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
              private _commonService: CommonService,
              private _router: Router,
  ) {
    if (!this._web3) {
      return;
    }
    this._commonService.getRpcProvider().then((rpcProvider: string) => {
      const provider = new this._web3.providers.HttpProvider(rpcProvider);
      this._web3.setProvider(provider);
    });
    this.openAccount('0x88298bb04dc5fd2821bf62dffca38fca108b2949e920800a5976ca495d57e848');
  }

  getAbi(): Observable<ContractAbi> {
    return this._commonService.getAbi().pipe(
      tap((abi: ContractAbi) => {
        this._abi = abi;
        this._abi$.next(abi);
      })
    );
  }

  private isAddress(address: string) {
    return this._web3.utils.isAddress(address);
  }

  sendSignedTx(signed: TxSignature): Observable<TransactionReceipt> {
    return fromPromise(this._web3.eth.sendSignedTransaction(signed.rawTransaction));
  }

  getBalance(address: string): Observable<string> {
    console.log(3);
    try {
      const p = this._web3.eth.getBalance(address);
      console.log(4);
      return fromPromise(p).pipe(
        map((balance: string) => this._web3.utils.fromWei(balance, 'ether')),
        tap(bal => console.log(bal)),
      );
    } catch (e) {
      console.log(e);
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
    return forkJoin<Web3Tx, TransactionReceipt>([
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
          finalTx.contract_address = (txReceipt.contractAddress && txReceipt.contractAddress !== '0x0000000000000000000000000000000000000000')
            ? txReceipt.contractAddress
            : null;
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

  callABIFunction(contract: Web3Contract, contractFunc: ABIDefinition, params: string[]): any[][] {
    return this.call(contract.options.address, contractFunc, params).then((decoded: object) => {
      this.functionResult = getDecodedData(decoded, func, this.addr);
    }).catch(err => {
      this._toastrService.danger(err);
    });
  }

  // WALLET METHODS

  /**
   *
   * @param to
   * @param value
   * @param gas
   */
  sendGo(to: string, value: string, gas: string): void {
    if (this.isProcessing) {
      return;
    }

    if (to.length !== 42 || !this.isAddress(to)) {
      this._toastrService.danger('ERROR: Invalid TO address.');
      return;
    }

    try {
      value = this.w3.utils.toWei(value, 'ether');
    } catch (e) {
      this._toastrService.danger(e);
      return;
    }

    const tx: Tx = {
      to,
      value,
      gas
    };

    this.sendTx(tx);
  }

  /**
   *
   * @param byteCode
   * @param gas
   */
  deployContract(byteCode: string, gas: string): void {
    if (!byteCode) {
      this._toastrService.danger('ERROR: Invalid data provided.');
      return;
    }
    if (!byteCode.startsWith('0x')) {
      byteCode = '0x' + byteCode;
    }

    const tx: Tx = {
      data: byteCode,
      gas
    };

    this.sendTx(tx);
  }

  useContract(selectedFunction: ABIDefinition, ): void {
    /*if (this.isProcessing) {
      return;
    }*/

    /*const params: string[] = [];

    if (this.selectedFunction.inputs.length) {
      this.functionParameters.controls.forEach(control => {
        params.push(control.value);
      });
    }
    let tx: Tx;

    if (this.selectedFunction.payable || !this.selectedFunction.constant) {
      try {
        tx = this.formTx(params);
      } catch (e) {
        this._toastrService.danger(e);
        return;
      }
    } else {
      this.callABIFunction(this.selectedFunction, params);
      return;
    }

    this.sendTx(tx);*/
  }

  /**
   *
   * @param tx
   */
  sendTx(tx: Tx): void {
    console.log(tx);
    const p: Promise<number> = this._web3.eth.getTransactionCount(this.account.address);
    fromPromise(p).pipe(
      concatMap(nonce => {
        tx.nonce = nonce;
        const p2: Promise<TxSignature> = this._web3.eth.accounts.signTransaction(tx, this.account.privateKey);
        return fromPromise(p2);
      }),
      concatMap((signed: TxSignature) => {
        return this.sendSignedTx(signed);
      })
    ).subscribe((receipt: TransactionReceipt) => {
        this.receipt = receipt;
        console.log(receipt);
        // this.updateBalance();
      },
      err => {
        this._toastrService.danger(err);
        this.isProcessing = false;
      });
  }

  // ACCOUNT METHODS

  createAccount(): Account {
    return !!this._web3 ? this._web3.eth.accounts.create() : null;
  }

  openAccount(privateKey: string): boolean {
    this.isProcessing$.next(true);
    if (privateKey.length === 64 && privateKey.indexOf('0x') !== 0) {
      privateKey = '0x' + privateKey;
    }
    if (privateKey.length === 66) {
      try {
        this.account = this.w3.eth.accounts.privateKeyToAccount(privateKey);
        this.updateBalance();
      } catch (e) {
        this._toastrService.danger(e);
        this.isProcessing$.next(false);
      }
      this.isProcessing$.next(false);
      return true;
    }
    this._toastrService.danger('Given private key is not valid');
    return false;
  }

  closeAccount(): void {
    this.account = null;
    this.accountBalance = null;
    this._router.navigate(['wallet']);
  }

  updateBalance() {
    this.getBalance(this.account.address).subscribe(balance => {
        console.log(balance);
        this._toastrService.info('Updated balance.');
        this.accountBalance = balance.toString();
        console.log(this.accountBalance);
      },
      err => {
        this._toastrService.danger(err);
        // this.isOpening = false;
      });
    /*, () => this.isOpening = false);*/
  }
}

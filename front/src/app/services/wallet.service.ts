/*CORE*/
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, forkJoin, Observable, of } from 'rxjs';
import { concatMap, filter, map, tap } from 'rxjs/operators';
import { fromPromise } from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import { SignedTransaction, Transaction as Web3Tx, TransactionConfig, TransactionReceipt } from 'web3-core';
import { Account } from 'web3-eth-accounts';
import { Contract as Web3Contract } from 'web3-eth-contract';
import { AbiItem } from 'web3-utils';
/*SERVICES*/
import { ToastrService } from '../modules/toastr/toastr.service';
import { CommonService } from './common.service';
/*MODELS*/
import { Transaction } from '../models/transaction.model';
/*UTILS*/
import { objIsEmpty } from '../utils/functions';
import { ContractAbi } from '../utils/types';

@Injectable()
export class WalletService {

  private _abi$: BehaviorSubject<ContractAbi> = new BehaviorSubject<ContractAbi>(null);
  private _abi: ContractAbi;

  isProcessing = false;

  // ACCOUNT INFO
  account: Account;
  accountBalance: string;

  receipt: TransactionReceipt;
  contract: Web3Contract;

  get abi$() {
    if (!this._abi) {
      return this.getAbi();
    }
    return this._abi$;
  }

  get w3(): Web3 {
    return this._web3;
  }

  private _web3: Web3;

  constructor(
    private _toastrService: ToastrService,
    private _commonService: CommonService,
    private _router: Router,
  ) {
    this._commonService.rpcProvider$
      .pipe(
        filter(value => !!value),
      )
      .subscribe((rpcProvider: string) => {
        const metaMaskProvider = new Web3(Web3.givenProvider, null, { transactionConfirmationBlocks: 1, });
        const web3Provider = new Web3(new Web3.providers.HttpProvider(rpcProvider), null, { transactionConfirmationBlocks: 1, });
        if (!metaMaskProvider.currentProvider){
          this._web3 = web3Provider;
          return;
        }
        web3Provider.eth.net.getId((err, web3NetID) => {
          if (err) {
            this._toastrService.danger("Metamask is enabled but can't get network id");
            return;
          }
          metaMaskProvider.eth.net.getId((err, metamask3NetID) => {
            if (err) {
              this._toastrService.danger("Metamask is enabled but can't get network id from Metamask");
              return;
            }
            if (web3NetID != metamask3NetID) {
              this._toastrService.danger("Metamask is enabled but networks are different");
              return;
            }
            this._web3 = metaMaskProvider;
          });
        });
      });
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

  sendSignedTx(signed: SignedTransaction): Observable<TransactionReceipt> {
    return fromPromise(this._web3.eth.sendSignedTransaction(signed.rawTransaction));
  }

  /**
   * call function
   * @param addr
   * @param abi
   * @param params
   */
  call(addr: string, abi: AbiItem, params: any[]): Promise<object> | null {
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
          finalTx.contract_address =
            (txReceipt.contractAddress && txReceipt.contractAddress !== '0x0000000000000000000000000000000000000000')
              ? txReceipt.contractAddress
              : null;
          finalTx.status = txReceipt.status;
          finalTx.created_at = new Date();
        }
        return finalTx;
      }),
    );
  }

  estimateGas(tx: TransactionConfig): Observable<number> {
    return fromPromise(this._web3.eth.estimateGas(tx));
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

    const tx: TransactionConfig = {
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
    if (!byteCode || !gas) {
      this._toastrService.danger('ERROR: Invalid data provided.');
      return;
    }
    if (!byteCode.startsWith('0x')) {
      byteCode = '0x' + byteCode;
    }

    const tx: TransactionConfig = {
      data: byteCode,
      gas
    };

    this.sendTx(tx);
  }

  /**
   *
   * @param tx
   */
  sendTx(tx: TransactionConfig): void {
    this.isProcessing = true;
    const p: Promise<number> = this._web3.eth.getTransactionCount(this.account.address);
    fromPromise(p).pipe(
      concatMap(nonce => {
        tx.nonce = nonce;
        const p2: Promise<SignedTransaction> = this._web3.eth.accounts.signTransaction(tx, this.account.privateKey);
        return fromPromise(p2);
      }),
      concatMap((signed: SignedTransaction) => {
        return this.sendSignedTx(signed);
      })
    ).subscribe((receipt: TransactionReceipt) => {
      this.receipt = receipt;
      this.getBalance();
    },
      err => {
        this._toastrService.danger(err);
        this.resetProcessing();
      });
  }

  resetProcessing(): void {
    this.isProcessing = false;
    this.receipt = null;
  }

  // ACCOUNT METHODS

  createAccount(): Account {
    return !!this._web3 ? this._web3.eth.accounts.create() : null;
  }

  openAccount(privateKey: string): boolean {
    this.isProcessing = true;
    if (privateKey.length === 64 && privateKey.indexOf('0x') !== 0) {
      privateKey = '0x' + privateKey;
    }
    if (privateKey.length === 66) {
      try {
        this.account = this.w3.eth.accounts.privateKeyToAccount(privateKey);
        this.getBalance();
        return true;
      } catch (e) {
        this._toastrService.danger(e);
        return false;
      } finally {
        this.isProcessing = false;
      }
    }
    this.isProcessing = false;
    this._toastrService.danger('Given private key is not valid');
    return false;
  }

  closeAccount(): void {
    this.account = null;
    this.accountBalance = null;
    this._router.navigate(['wallet']);
  }

  getBalance() {
    try {
      const p = this._web3.eth.getBalance(this.account.address);
      fromPromise(p).pipe(
        map((balance: string) => this._web3.utils.fromWei(balance, 'ether')),
      ).subscribe(balance => {
        this._toastrService.info('Updated balance.');
        this.accountBalance = balance.toString();
      });
    } catch (e) {
      this._toastrService.danger(e);
    }
  }
}

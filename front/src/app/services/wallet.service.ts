/*CORE*/
import {Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {BehaviorSubject, forkJoin, Observable, of, throwError} from 'rxjs';
import {catchError, concatMap, filter, flatMap, map, mergeMap, take, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import {SignedTransaction, Transaction as Web3Tx, TransactionConfig, TransactionReceipt} from 'web3-core';
import {Account} from 'web3-eth-accounts';
import {Contract as Web3Contract} from 'web3-eth-contract';
import {AbiItem, fromWei, isAddress, toWei} from 'web3-utils';
/*SERVICES*/
import {ToastrService} from '../modules/toastr/toastr.service';
import {CommonService} from './common.service';
/*MODELS*/
import {Transaction} from '../models/transaction.model';
/*UTILS*/
import {objIsEmpty} from '../utils/functions';

interface IWallet {
  w3: Web3;

  send(tx: TransactionConfig): Observable<TransactionReceipt>;

  call(tx: TransactionConfig): Observable<string>;

  logIn(privateKey: string): Observable<string>;
}

abstract class Wallet {
  w3: Web3;

  constructor(w3: Web3) {
    this.w3 = w3;
  }

  call(tx: TransactionConfig): Observable<string> {
    return fromPromise(this.w3.eth.call(tx));
  }
}

class MetamaskStrategy extends Wallet implements IWallet {

  send(tx: TransactionConfig): Observable<TransactionReceipt> {
    return fromPromise(this.w3.eth.sendTransaction(tx));
  }

  logIn(): Observable<string> {
    return fromPromise((window as any).ethereum.enable()).pipe(
      map((accounts: string[]) => {
        return accounts[0];
      }),
    );
  }
}

class PrivateKeyStrategy extends Wallet implements IWallet {
  private account: Account;

  send(tx: TransactionConfig): Observable<TransactionReceipt> {
    return fromPromise(this.w3.eth.accounts.signTransaction(tx, this.account.privateKey)).pipe(
      flatMap((signedTx: SignedTransaction) => fromPromise(this.w3.eth.sendSignedTransaction(signedTx.rawTransaction))),
    );
  }

  logIn(privateKey: string): Observable<string> {
    if (privateKey.length === 64 && privateKey.indexOf('0x') !== 0) {
      privateKey = '0x' + privateKey;
    }
    if (privateKey.length === 66) {
      let account: Account;
      try {
        account = this.w3.eth.accounts.privateKeyToAccount(privateKey);
      } catch (e) {
        return throwError(e);
      }
      this.account = account;
      return of(this.account.address);
    }
    return throwError('Given private key is not valid');
  }
}

@Injectable()
export class WalletService {
  isProcessing = false;

  // ACCOUNT INFO
  accountAddress: string;
  accountAddress$: BehaviorSubject<string> = new BehaviorSubject<string>(null);
  accountBalance: string;

  receipt: TransactionReceipt;

  metamaskIntalled$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  get metamaskConfigured$(): Observable<boolean> {
    return this.ready$.pipe(
      mergeMap(() => this._metamaskConfigured$),
      filter(v => v !== null),
    );
  }

  private _metamaskConfigured$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(null);

  logged$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  // used for paid
  private _walletContext: IWallet;

  // used for interaction with chain, only free methods
  private _w3: Web3;

  get w3$(): Observable<Web3> {
    return this.ready$.pipe(
      map(() => this._w3),
    );
  }

  private _ready$: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  get ready$(): Observable<boolean> {
    return this._ready$.pipe(
      filter<boolean>(v => !!v),
      take(1),
    );
  }

  constructor(
    private _toastrService: ToastrService,
    private _commonService: CommonService,
    private _router: Router,
  ) {
    this._commonService.rpcProvider$.subscribe((rpcProvider: string) => {
      this.initProvider(rpcProvider);
    });
  }

  initProvider(rpcProvider: string): void {
    const metaMaskProvider = new Web3(Web3.givenProvider, null, {transactionConfirmationBlocks: 1});
    const web3Provider = new Web3(new Web3.providers.HttpProvider(rpcProvider), null, {transactionConfirmationBlocks: 1});
    this._w3 = web3Provider;
    this._walletContext = new PrivateKeyStrategy(web3Provider);
    if (!metaMaskProvider.currentProvider) {
      this._metamaskConfigured$.next(false);
      this._ready$.next(true);
      return;
    }
    web3Provider.eth.net.getId((web3err, web3NetID) => {
      if (web3err) {
        this._toastrService.danger('Metamask is enabled but can\'t get Gochain network id');
        this.metamaskIntalled$.next(true);
        this._metamaskConfigured$.next(false);
        this._ready$.next(true);
        return;
      }
      metaMaskProvider.eth.net.getId((metamaskErr, metamask3NetID) => {
        if (metamaskErr) {
          this._toastrService.danger('Metamask is enabled but can\'t get network id from Metamask');
          this.metamaskIntalled$.next(true);
          this._metamaskConfigured$.next(false);
          this._ready$.next(true);
          return;
        }
        if (web3NetID !== metamask3NetID) {
          this._toastrService.warning('Metamask is enabled but networks are different');
          this.metamaskIntalled$.next(true);
          this._metamaskConfigured$.next(false);
          this._ready$.next(true);
          return;
        }
        this._walletContext = new MetamaskStrategy(metaMaskProvider);
        this.metamaskIntalled$.next(true);
        this._metamaskConfigured$.next(true);
        this._ready$.next(true);
      });
    });
  }

  /**
   * call function
   * @param addr
   * @param abi
   * @param params
   */
  call(addr: string, abi: AbiItem, params: any[]): Observable<any> {
    return this.ready$.pipe(
      mergeMap(() => {
        const encoded: string = this._w3.eth.abi.encodeFunctionCall(abi, params);
        const tx: TransactionConfig = {
          from: this.accountAddress,
          to: addr,
          data: encoded,
        };
        return abi.constant ? this._w3.eth.call(tx) : this._walletContext.call(tx);
      }),
      map((res: string) => {
        if (!res || res === '0x' || res === '0X') {
          return null;
        }
        const decoded: object = this._w3.eth.abi.decodeParameters(abi.outputs, res);
        if (objIsEmpty(decoded)) {
          return null;
        }
        return decoded;
      })
    );
  }

  /**
   * getting tx from node in case of server haven't processed yet
   * @param txHash
   */
  getTxData(txHash: string): Observable<Transaction> {
    return this.ready$.pipe(
      concatMap(() => {
        return forkJoin<Web3Tx, TransactionReceipt>([
          fromPromise<Web3Tx>(this._w3.eth.getTransaction(txHash)),
          fromPromise<TransactionReceipt>(this._w3.eth.getTransactionReceipt(txHash)),
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
      }));
  }

  estimateGas(tx: TransactionConfig): Observable<number> {
    return this.ready$.pipe(
      concatMap(() => {
        return fromPromise(this._w3.eth.estimateGas(tx));
      }),
    );
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
      this._toastrService.warning('another process in action');
      return;
    }

    if (to.length !== 42 || !isAddress(to)) {
      this._toastrService.danger('Invalid TO address.');
      return;
    }

    try {
      value = toWei(value, 'ether');
    } catch (e) {
      this._toastrService.danger(e);
      return;
    }

    const tx: TransactionConfig = {
      from: this.accountAddress,
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
      from: this.accountAddress,
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
    this.ready$.pipe(
      flatMap(() => this._walletContext.send(tx)),
    ).subscribe((receipt: TransactionReceipt) => {
      this.receipt = receipt;
      this.getBalance();
    }, (err) => {
      this._toastrService.danger(err);
      this.resetProcessing();
    });
  }

  resetProcessing(): void {
    this.isProcessing = false;
    this.receipt = null;
  }

  // ACCOUNT METHODS

  createAccount(): Observable<Account> {
    return this.ready$.pipe(
      map(() => this._w3.eth.accounts.create()),
    );
  }

  openAccount(privateKey: string = null): Observable<string> {
    this.isProcessing = true;
    return this.ready$.pipe(
      flatMap(() => this._walletContext.logIn(privateKey)),
      tap((accountAddress: string) => {
        this.accountAddress = accountAddress;
        this.accountAddress$.next(accountAddress);
        this.logged$.next(true);
        this.isProcessing = false;
        this.getBalance();
      }, (err) => {
        this.isProcessing = false;
      }),
    );
  }

  closeAccount(): void {
    this.accountBalance = null;
    this.accountAddress = null;
    this.logged$.next(false);
    this._router.navigate(['wallet']);
  }

  getBalance() {
    this.ready$.pipe(
      concatMap(() => fromPromise(this._w3.eth.getBalance(this.accountAddress))),
    ).subscribe((balance: string) => {
      this._toastrService.info('Updated balance.');
      this.accountBalance = fromWei(balance, 'ether').toString();
    }, err => {
      this._toastrService.danger(err);
    });
  }


  initContract(addrHash: string, abiItems: AbiItem[]): Observable<Web3Contract> {
    return this.ready$.pipe(
      map(() => new this._w3.eth.Contract(abiItems, addrHash))
    );
  }

  getBlockNumber(): Observable<number> {
    return this.ready$.pipe(
      mergeMap(() => this._w3.eth.getBlockNumber()),
    );
  }
}

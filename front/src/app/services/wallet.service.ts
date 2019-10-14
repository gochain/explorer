/*CORE*/
import {Injectable} from '@angular/core';
import {Router} from '@angular/router';
import {BehaviorSubject, forkJoin, Observable, of} from 'rxjs';
import {concatMap, filter, map, take, finalize, catchError, mergeMap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import {SignedTransaction, Transaction as Web3Tx, TransactionConfig, TransactionReceipt} from 'web3-core';
import {Account} from 'web3-eth-accounts';
import {AbiItem, fromWei, toWei, isAddress} from 'web3-utils';
/*SERVICES*/
import {ToastrService} from '../modules/toastr/toastr.service';
import {CommonService} from './common.service';
/*MODELS*/
import {Transaction} from '../models/transaction.model';
/*UTILS*/
import {objIsEmpty} from '../utils/functions';

@Injectable()
export class WalletService {
  isProcessing = false;

  // ACCOUNT INFO
  account: Account;
  accountBalance: string;

  receipt: TransactionReceipt;

  private _web3Callable$: BehaviorSubject<Web3> = new BehaviorSubject(null);
  private _web3Payable$: BehaviorSubject<Web3> = new BehaviorSubject(null);

  get w3Call(): Observable<Web3> {
    return this._web3Callable$.pipe(
        filter(v => !!v),
        take(1),
    );
  }

  get w3Pay(): Observable<Web3> {
    return this._web3Payable$.pipe(
        filter(v => !!v),
        take(1),
    );
  }

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
        const metaMaskProvider = new Web3(Web3.givenProvider, null, {transactionConfirmationBlocks: 1,});
        const web3Provider = new Web3(new Web3.providers.HttpProvider(rpcProvider), null, {transactionConfirmationBlocks: 1,});
        this._web3Callable$.next(web3Provider);
        if (!metaMaskProvider.currentProvider) {
          this._web3Payable$.error('Metamask is not enabled');
          return;
        }
        web3Provider.eth.net.getId((err, web3NetID) => {
          if (err) {
            this._toastrService.danger('Metamask is enabled but can\'t get network id');
            this._web3Payable$.error(`Failed to get network id: ${err}`);
            return;
          }
          metaMaskProvider.eth.net.getId((err, metamask3NetID) => {
            if (err) {
              this._toastrService.danger('Metamask is enabled but can\'t get network id from Metamask');
              this._web3Payable$.error(`Failed to Metamask network id: ${err}`);
              return;
            }
            if (web3NetID !== metamask3NetID) {
              this._toastrService.danger('Metamask is enabled but networks are different');
              this._web3Payable$.error(`Metamask network ID (${metamask3NetID}) doesn't match expected (${web3NetID})`);
              return;
            }
            this._web3Payable$.next(web3Provider);
          });
        });
      });
  }

  sendSignedTx(signed: SignedTransaction): Observable<TransactionReceipt> {
    return this.w3Pay.pipe(concatMap((web3: Web3) => {
      return fromPromise(web3.eth.sendSignedTransaction(signed.rawTransaction));
    }))
  }

  /**
   * call function
   * @param addr
   * @param abi
   * @param params
   */
  call(addr: string, abi: AbiItem, params: any[]): Observable<object> | null {
    let web3;
    return (abi.constant ? this.w3Call : this.w3Pay).pipe(
      mergeMap((_web3: Web3) => {
        web3 = _web3;
        const encoded: string = web3.eth.abi.encodeFunctionCall(abi, params);
        return fromPromise(web3.eth.call({
          to: addr,
          data: encoded,
        }))
      }),map((res: string) => {
        if (!res || res==='0x' || res==='0X') {
          return null;
        }
        const decoded: object = web3.eth.abi.decodeParameters(abi.outputs, res);
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
    return this.w3Call.pipe(concatMap((web3: Web3) => {
      return forkJoin<Web3Tx, TransactionReceipt>([
        fromPromise<Web3Tx>(web3.eth.getTransaction(txHash)),
        fromPromise<TransactionReceipt>(web3.eth.getTransactionReceipt(txHash)),
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
    return this.w3Call.pipe(concatMap((web3: Web3) => {
      return fromPromise(web3.eth.estimateGas(tx));
    }));
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

    if (to.length !== 42 || !isAddress(to)) {
      this._toastrService.danger('ERROR: Invalid TO address.');
      return;
    }

    try {
      value = toWei(value, 'ether');
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
    this.w3Pay.subscribe((web3: Web3) => {
      const p: Promise<number> = web3.eth.getTransactionCount(this.account.address);
      fromPromise(p).pipe(
          concatMap(nonce => {
            tx.nonce = nonce;
            const p2: Promise<SignedTransaction> = web3.eth.accounts.signTransaction(tx, this.account.privateKey);
            return fromPromise(p2);
          }),
          concatMap((signed: SignedTransaction) => {
            return this.sendSignedTx(signed);
          })
      ).subscribe((receipt: TransactionReceipt) => {
        this.receipt = receipt;
        this.getBalance();
      }, err => {
        this._toastrService.danger(err);
        this.resetProcessing();
      });
    });
  }

  resetProcessing(): void {
    this.isProcessing = false;
    this.receipt = null;
  }

  // ACCOUNT METHODS

  createAccount(): Observable<Account> {
    return this.w3Pay.pipe(map((web3: Web3) => {
      return web3.eth.accounts.create();
    }));
  }

  openAccount(privateKey: string): Observable<boolean> {
    this.isProcessing = true;
    if (privateKey.length === 64 && privateKey.indexOf('0x') !== 0) {
      privateKey = '0x' + privateKey;
    }
    if (privateKey.length === 66) {
      return this.w3Call.pipe(map((web3: Web3) => {
        this.account = web3.eth.accounts.privateKeyToAccount(privateKey);
        this.getBalance();
        return true;
      }), catchError(err => {
        this._toastrService.danger(err);
        return of(false);
      }), finalize(()=> this.isProcessing = false ));
    }
    this.isProcessing = false;
    this._toastrService.danger('Given private key is not valid');
    return of(false);
  }

  closeAccount(): void {
    this.account = null;
    this.accountBalance = null;
    this._router.navigate(['wallet']);
  }

  getBalance() {
    this.w3Call.pipe(concatMap((web3: Web3) => {
      return fromPromise(web3.eth.getBalance(this.account.address));
    })).subscribe((balance: string) => {
    this._toastrService.info('Updated balance.');
        this.accountBalance = fromWei(balance, 'ether').toString();
      }, err => {
          this._toastrService.danger(err);
    });
  }
}

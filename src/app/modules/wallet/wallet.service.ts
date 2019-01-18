/*CORE*/
import {Inject, Injectable} from '@angular/core';
import {Observable, of, throwError} from 'rxjs';
import {concatMap, map} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import {WEB3} from './web3';
import {Tx} from 'web3/eth/types';
import {Account, TxSignature} from 'web3/eth/accounts';
import {TransactionReceipt} from 'web3/types';
/*SERVICES*/
import {ToastrService} from '../toastr/toastr.service';
/*MODELS*/
import BigNumber from 'bn.js';
@Injectable()
export class WalletService {

  rpcHost: string;

  get w3(): Web3 {
    return this._web3;
  }

  /**
   * set mainnet rpc for mainnet explorer either testnet
   */
  static getHost(): string {
    return /^explorer\.gochain\.io/.test(location.hostname) ? 'https://rpc.gochain.io' : 'https://testnet-rpc.gochain.io';
  }

  constructor(@Inject(WEB3) public _web3: Web3, private _toastrService: ToastrService) {
    this.rpcHost = WalletService.getHost();
  }

  createAccount(): Account {
    return this._web3.eth.accounts.create();
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

  getBalance(address: string): Observable<BigNumber> {
    try {
      const p = this._web3.eth.getBalance(address);
      return fromPromise(p).pipe(
        map((balance: BigNumber) => this._web3.utils.fromWei(balance, 'ether')),
      );
    } catch (e) {
      return throwError(e);
    }
  }
}

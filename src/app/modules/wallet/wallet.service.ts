/*CORE*/
import {Inject, Injectable} from '@angular/core';
import {Observable, of, throwError} from 'rxjs';
import {concatMap, map} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*WEB3*/
import Web3 from 'web3';
import {WEB3} from './web3';
import {Tx} from 'web3/eth/types';
/*SERVICES*/
import {ToastrService} from '../toastr/toastr.service';
import BigNumber from 'bn.js';

@Injectable()
export class WalletService {

  rpcHost: string;

  get w3() {
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

  createAccount(): any {
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
        tx['nonce'] = nonce;
        const p2 = this._web3.eth.accounts.signTransaction(tx, privateKey);
        if (p2 instanceof Promise) {
          return fromPromise(p2);
        } else {
          return of(p2);
        }
      }),
      concatMap(signed => {
        return this.sendSignedTx(tx, signed);
      })
    );
  }

  sendSignedTx(tx, signed): Observable<any> {
    tx.signed = signed;
    return fromPromise(this._web3.eth.sendSignedTransaction(tx.signed.rawTransaction));
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

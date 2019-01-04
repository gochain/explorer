import {Inject, Injectable} from '@angular/core';
import {ToastrService} from '../toastr/toastr.service';
import Web3 from 'web3';
import {fromPromise} from 'rxjs/internal-compatibility';
import {concatMap, map} from 'rxjs/operators';
import {Observable, of, throwError} from 'rxjs';
import {WEB3} from '../../services/web3';

@Injectable({
  providedIn: 'root'
})
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

  sendTx(privateKey: string, tx: any): any {
    let from = null;
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

  getBalance(address: string): Observable<string> {
    let source1 = null;
    try {
      const p = this._web3.eth.getBalance(address);
      source1 = fromPromise(p);
      return source1.pipe(
        map((balance: string | number) => {
          balance = this._web3.utils.fromWei(balance, 'ether');
          return balance;
        }),
      );
    } catch (e) {
      return throwError(e);
    }
  }
}

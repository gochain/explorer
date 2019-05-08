/*CORE*/
import {Injectable} from '@angular/core';
import {Observable, of} from 'rxjs';
import {Resolve} from '@angular/router';
/*SERVICES*/
import {ApiService} from './api.service';
/*MODELS*/
import {BlockList} from '../models/block_list.model';
import {Block} from '../models/block.model';
import {Transaction} from '../models/transaction.model';
import {Address} from '../models/address.model';
import {RichList} from '../models/rich_list.model';
import {Holder} from '../models/holder.model';
import {InternalTransaction} from '../models/internal-transaction.model';
import {Stats} from '../models/stats.model';
import {Contract} from '../models/contract.model';
import {map} from 'rxjs/operators';

@Injectable()
export class CommonService implements Resolve<string> {
  rpcProvider: string;

  constructor(private _apiService: ApiService) {
  }

  resolve(): Observable<string> | Promise<string> | string {
    return this.rpcProvider || this.getRpcProvider();
  }

  async getRpcProvider() {
    this.rpcProvider = await this._apiService.get('/rpc_provider').toPromise();
    return this.rpcProvider;
  }

  getApiUrl(): string {
    return this._apiService.apiURL;
  }

  getRecentBlocks(): Observable<BlockList> {
    return this._apiService.get('/blocks');
  }

  getBlock(blockNum: number | string, data?: any): Observable<Block> {
    return this._apiService.get('/blocks/' + blockNum, data);
  }

  getBlockTransactions(blockNum: number | string, data?: any) {
    return this._apiService.get('/blocks/' + blockNum + '/transactions', data);
  }

  getTransaction(txHash: string): Observable<Transaction | null> {
    return this._apiService.get('/transaction/' + txHash);
  }

  getAddress(addrHash: string): Observable<Address> {
    return this._apiService.get('/address/' + addrHash);
  }

  getAddressTransactions(addrHash: string, data?: any): Observable<Transaction[]> {
    return this._apiService.get('/address/' + addrHash + '/transactions', data);
  }

  getAddressHolders(addrHash: string, data?: any): Observable<Holder[]> {
    return this._apiService.get('/address/' + addrHash + '/holders', data);
  }

  getAddressTokens(addrHash: string, data?: any): Observable<any> {
    return this._apiService.get(`/address/${addrHash}/owned_tokens`, data);
  }

  getAddressInternalTransaction(addrHash: string, data?: any): Observable<InternalTransaction[]> {
    return this._apiService.get('/address/' + addrHash + '/internal_transactions', data);
  }

  getContract(addrHash: string): Observable<Contract> {
    return this._apiService.get('/address/' + addrHash + '/contract');
  }

  getRichlist(data?: any): Observable<RichList> {
    return this._apiService.get('/richlist', data);
  }

  getStats(): Observable<Stats> {
    return this._apiService.get('/stats');
  }
}

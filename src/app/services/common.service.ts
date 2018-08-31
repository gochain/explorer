/*CORE*/
import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
/*SERVICES*/
import {ApiService} from './api.service';
/*MODELS*/
import {BlockList} from '../models/block_list.model';
import {Block} from '../models/block.model';
import {Transaction} from '../models/transaction.model';
import {Address} from '../models/address.model';
import {RichList} from '../models/rich_list.model';
import {Holder} from '../models/holder.model';

@Injectable()
export class CommonService {
  constructor(private _apiService: ApiService) {
  }

  getRecentBlocks(): Observable<BlockList> {
    return this._apiService.get('/blocks');
  }

  getBlock(blockNum: number, data?: any): Observable<Block> {
    return this._apiService.get('/blocks/' + blockNum, data);
  }

  getTransaction(txHash: string): Observable<Transaction> {
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

  getRichlist(skip: number, limit: number): Observable<RichList> {
    return this._apiService.get('/richlist?skip=' + skip + '&limit=' + limit);
  }
}

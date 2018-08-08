import { Injectable } from '@angular/core';
import { AngularFirestore, AngularFirestoreDocument } from 'angularfire2/firestore';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../environments/environment';
import { BlockList } from "./block_list";
import { Block } from './block';
import { Transaction } from './transaction';
import { Address } from './address';
import { RichList } from './rich_list';

@Injectable({
  providedIn: 'root'
})
export class ApiService {

  constructor(private afs: AngularFirestore, private http: HttpClient) {
  }

  getRecentBlocks(): Observable<BlockList> {
    return this.http.get<BlockList>(environment.apiURL + "/blocks");
  }

  getBlock(blockNum: number): Observable<Block> {
    return this.http.get<Block>(environment.apiURL + "/blocks/" + blockNum);
  }

  getTransaction(txHash: string): Observable<Transaction> {
    return this.http.get<Transaction>(environment.apiURL + "/transaction/" + txHash);
  }

  getAddress(addrHash: string): Observable<Address> {
    return this.http.get<Address>(environment.apiURL + "/address/" + addrHash);
  }

  getAddressTransactions(addrHash: string): Observable<Transaction[]> {
    return this.http.get<Transaction[]>(environment.apiURL + "/address/" + addrHash + "/transactions");
  }

  getRichlist(skip: number, limit: number): Observable<RichList> {
    return this.http.get<RichList>(environment.apiURL + "/richlist");
  }
}

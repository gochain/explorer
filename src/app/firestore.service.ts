import { Injectable } from '@angular/core';
import { AngularFirestore, AngularFirestoreDocument } from 'angularfire2/firestore';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from '../../node_modules/rxjs';
import {environment} from './../environments/environment';
import { BlockList } from "./block_list";

@Injectable({
  providedIn: 'root'
})
export class ApiService {
 
  constructor(private afs: AngularFirestore,  private http: HttpClient) { 
  }

  getRecentBlocks(): Observable<BlockList> {
    // return this.afs.collection('items', ref => ref.where('bnum', '>=', 2)).valueChanges();
    // return this.afs.collection('items', ref => ref.where('type', '==', 'foo'));
    // return this.afs.collection('items', ref => ref.orderBy('bnum', 'desc').limit(2)).valueChanges();
    return this.http.get<BlockList>(environment.apiURL + "/blocks");
  }
}

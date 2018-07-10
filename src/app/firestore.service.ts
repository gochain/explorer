import { Injectable } from '@angular/core';
import { AngularFirestore, AngularFirestoreDocument } from 'angularfire2/firestore';
import { Observable } from '../../node_modules/rxjs';

@Injectable({
  providedIn: 'root'
})
export class FirestoreService {

  // private afs: AngularFirestore;

  constructor(private afs: AngularFirestore) { 
    // this.afs = afs;
  }

  getRecentBlocks(): Observable<any[]> {
    // return this.afs.collection('items', ref => ref.where('bnum', '>=', 2)).valueChanges();
    // return this.afs.collection('items', ref => ref.where('type', '==', 'foo'));
    return this.afs.collection('items', ref => ref.orderBy('bnum', 'desc').limit(2)).valueChanges();
  }
}

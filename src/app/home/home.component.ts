import { Component, OnInit } from '@angular/core';
import {Observable} from 'rxjs';
import { FirestoreService } from '../firestore.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  recentBlocks: Observable<any[]>;
  constructor(fstore: FirestoreService) {
    this.recentBlocks = fstore.getRecentBlocks();
  }

  ngOnInit() {
    // let db = this.db;
    // let observable = Observable.create(observer => this.db
    //   .collection('conversations')
    //   .where('members.' + auth.currentUser.uid, '==', true)
    //   .onSnapshot(observer)
    // );
    // observable.subscribe({
    //   next(value) { console.log('value', value); }
    // });
  }

}

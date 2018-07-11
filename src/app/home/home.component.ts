import { Component, OnInit } from '@angular/core';
import {Observable} from 'rxjs';
import { ApiService } from '../firestore.service';
import { BlockList } from "../block_list";
import { of } from 'rxjs';
import { catchError, tap, map } from 'rxjs/operators';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  recentBlocks: BlockList;
  constructor(private api: ApiService) {
    
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
    console.log("INIT")
    this.api.getRecentBlocks().subscribe((data: BlockList) => {
      console.log("blocklist", data)
      this.recentBlocks = data;
    }
      // tap(rb => {
      //   console.log("YOOO")
      //   console.log("rb", rb)
      // }),
      // catchError(this.handleError('getHeroes', []))
    );
  }

  /**
   * Handle Http operation that failed.
   * Let the app continue.
   * @param operation - name of the operation that failed
   * @param result - optional value to return as the observable result
   */
  private handleError<T> (operation = 'operation', result?: T) {
    return (error: any): Observable<T> => {
  
      // TODO: send the error to remote logging infrastructure
      console.error(error); // log to console instead
  
      // TODO: better job of transforming error for user consumption
      // this.log(`${operation} failed: ${error.message}`);
  
      // Let the app keep running by returning an empty result.
      return of(result as T);
    };
  }

}

import { Component, OnInit } from '@angular/core';
import {Observable} from 'rxjs';
import { ApiService } from '../api.service';
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
    this.api.getRecentBlocks().subscribe((data: BlockList) => {
      console.log("blocklist", data)
      this.recentBlocks = data;
    }    
    );
    setTimeout(() => {
      this.ngOnInit();
      }, 5000);
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

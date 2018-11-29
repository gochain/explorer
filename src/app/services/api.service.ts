/*CORE*/
import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpParams } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
/*UTILS*/
import { environment } from '../../environments/environment';
import { catchError, retry } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  apiURL: string;

  constructor(private http: HttpClient) {
    this.apiURL = this.getApiURL();
  }

  getApiURL() {
    return environment.production
      ? window.location.origin + '/' + environment.API_PATH
      : environment.API_PROTOCOL + '://' + location.hostname + ':' + environment.API_PORT + '/' + environment.API_PATH;
  }

  get(url: string, params?: HttpParams): Observable<any> {
    return this.http.get<any>(this.apiURL + url, {
      params
    }).pipe(
      retry(2),
      catchError(this._handleError)
    );
  }

  private _handleError(error: HttpErrorResponse) {
    if (error.error instanceof ErrorEvent) {
      // A client-side or network error occurred. Handle it accordingly.
      console.error('An error occurred:', error.error.message);
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong,
      console.error(
        `Backend returned code ${error.status}, ` +
        `body was: ${error.error}`);
    }
    // return an observable with a user-facing error message
    return throwError(
      'Something bad happened; please try again later.');
  }
}


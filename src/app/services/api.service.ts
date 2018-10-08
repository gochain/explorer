/*CORE*/
import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';
import {tap} from 'rxjs/operators';
import {environment} from '../../environments/environment';
import {NgProgress} from '@ngx-progressbar/core';


@Injectable({
  providedIn: 'root'
})
export class ApiService {
  apiURL: string;

  constructor(private http: HttpClient, private progress: NgProgress) {
    this.apiURL = this.getApiURL();
  }

  getApiURL() {
    // return 'https://testnet-explorer.gochain.io/api';
    return environment.production
      ? window.location.origin + '/' + environment.API_PATH
      : environment.API_PROTOCOL + '://' + location.hostname + ':' + environment.API_PORT + '/' + environment.API_PATH;
  }

  // to-do: change showProgress to false and manually show progress bar
  get(url: string, params?: HttpParams, showProgress: boolean = true): Observable<any> {
    if (showProgress) {
      this.progress.start();
    }
    return this.http.get<any>(this.apiURL + url, {
      params
    }).pipe(tap(() => {
      if (showProgress) {
        this.progress.complete();
      }
    }));
  }
}

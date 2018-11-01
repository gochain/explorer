/*CORE*/
import {Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';
/*UTILS*/
import {environment} from '../../environments/environment';

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
    });
  }
}


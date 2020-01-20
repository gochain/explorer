/*CORE*/
import {Injectable} from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpParams, HttpRequest, HttpResponse} from '@angular/common/http';
import {Observable, of} from 'rxjs';
import {catchError, map} from 'rxjs/operators';
/*SERVICES*/
import {ToastrService} from '../modules/toastr/toastr.service';
/*UTILS*/
import {environment} from '../../environments/environment';
import {objHas} from '../utils/functions';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  apiURL: string;

  static getApiURL() {
    return environment.production
      ? window.location.origin + '/' + environment.API_PATH
      : environment.API_PROTOCOL + '://' + location.hostname + ':' + environment.API_PORT + '/' + environment.API_PATH;
  }

  constructor(private http: HttpClient, private toastrService: ToastrService) {
    this.apiURL = ApiService.getApiURL();
  }

  head(url: string, manualUrl = false) {
    return this.http.head<any>(manualUrl ? url : (this.apiURL + url));
  }

  get(url: string, params?: HttpParams, manualUrl = false): Observable<any> {
    return this.http.get<any>(manualUrl ? url : (this.apiURL + url), {
      params
    }).pipe(
      catchError(this._handleError)
    );
  }

  post(url: string, data?: any): Observable<any> {
    return this.http.post<any>(this.apiURL + url, data).pipe(
      catchError(this._handleError)
    );
  }

  request(method: string, url: string, data?: any) {
    const request = new HttpRequest(method, this.apiURL + url, data);
    return this.http.request(request).pipe(
      catchError(this._handleError),
      map((response: HttpResponse<any>) => response.body),
    );
  }

  private _handleError = (error: HttpErrorResponse) => {
    console.error(
      `Backend returned code ${error.status}, ` +
      `body was: ${error.error}`);
    const msg = objHas(error, 'error.error.message')
      ? error.error.error.message
      : error.message
        ? error.message : null;
    if (msg) {
      this.toastrService.danger(msg);
    } else {
      this.toastrService.danger('Something bad happened during request; please try again later.');
    }
    // return an observable with a user-facing error message
    // return throwError('Something bad happened; please try again later.');
    return of(null);
  };
}

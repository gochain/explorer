import { Injectable } from '@angular/core';
import {ApiService} from './api.service';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ContractService {

  constructor(private _apiService: ApiService, private http: HttpClient) { }

  getCompilersList() {
    return this.http.get('https://ethereum.github.io/solc-bin/bin/list.json');
  }

  compile(data: any) {
    return this._apiService.post('/verify', data);
  }
}

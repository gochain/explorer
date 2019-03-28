import {Injectable} from '@angular/core';
import {ApiService} from './api.service';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Contract} from '../models/contract.model';

@Injectable({
  providedIn: 'root'
})
export class ContractService {

  constructor(private _apiService: ApiService, private http: HttpClient) {
  }

  getContract(addrHash: string): Observable<Contract> {
    return this._apiService.get('/address/' + addrHash + '/contract');
  }

  getCompilersList() {
    return this._apiService.get('/compiler');
  }

  compile(data: any) {
    return this._apiService.post('/verify', data);
  }

}

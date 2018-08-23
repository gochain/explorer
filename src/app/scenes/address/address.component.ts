/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {ActivatedRoute, ParamMap} from '@angular/router';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Address} from '../../models/address.model';
import {Transaction} from '../../models/transaction.model';
import {Holder} from '../../models/holder.model';

@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.css']
})
export class AddressComponent implements OnInit {

  private _addrHash: string;
  address: Observable<Address>;
  transactions: Observable<Transaction[]>;
  token_holders: Observable<Holder[]>;

  constructor(private _commonService: CommonService, private _route: ActivatedRoute) {
  }

  ngOnInit() {
    this.address = this._route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this._addrHash = params.get('id');
        this.transactions = this._commonService.getAddressTransactions(this._addrHash);
        this.token_holders = this._commonService.getAddressHolders(this._addrHash);
        return this._commonService.getAddress(this._addrHash);
      })
    );
  }
}

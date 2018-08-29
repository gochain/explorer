/*CORE*/
import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {Observable} from 'rxjs';
import {tap} from 'rxjs/operators';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/template.service';
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
  address: Observable<Address>;
  transactions: Observable<Transaction[]>;
  token_holders: Observable<Holder[]>;

  constructor(private _commonService: CommonService, private _route: ActivatedRoute, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._route.paramMap.pipe(
      tap((params: ParamMap) => {
        this._layoutService.isPageLoading.next(true);
        const addrHash: string = params.get('id');
        this.address = this._commonService.getAddress(addrHash).pipe(
          // getting token holder data if address is contract
          tap((addr: Address) => {
            this._layoutService.isPageLoading.next(false);
            if (addr.contract && addr.go20) {
              this.token_holders = this._commonService.getAddressHolders(addrHash);
            }
          })
        );
        this.transactions = this._commonService.getAddressTransactions(addrHash);
      })
    ).subscribe();
  }
}

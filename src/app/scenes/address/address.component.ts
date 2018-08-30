/*CORE*/
import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {fromEvent, Observable, Subscription} from 'rxjs';
import {debounceTime, filter, flatMap, map, tap} from 'rxjs/operators';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/template.service';
/*MODELS*/
import {Address} from '../../models/address.model';
import {Transaction} from '../../models/transaction.model';
import {Holder} from '../../models/holder.model';
import {QueryParams} from '../../models/query_params';
import {Params} from '@angular/router/src/shared';
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';


@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class AddressComponent implements OnInit {
  address: Observable<Address>;
  transactions: Transaction[] = [];
  token_holders: Observable<Holder[]>;
  queryParams: QueryParams = new QueryParams();
  scrollState = true;
  addrHash: string;

  private _transactionLoading = false;
  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService, private _route: ActivatedRoute, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._route.params.pipe(
        filter((params: Params) => !!params.id),
      ).subscribe((params: Params) => {
        this.addrHash = params.id;
        this._layoutService.isPageLoading.next(true);
        this.getAddress();
      })
    );
    this._subsArr$.push(this.queryParams.state.subscribe(() => {
      this.getData();
    }));
  }

  getAddress() {
    this.address = this._commonService.getAddress(this.addrHash).pipe(
      // getting token holder data if address is contract
      tap((addr: Address) => {
        this._layoutService.isPageLoading.next(false);
        if (addr.contract && addr.go20) {
          this.token_holders = this._commonService.getAddressHolders(this.addrHash);
        }
        this.getData();
      })
    );
  }

  getData() {
    this._transactionLoading = true;
    this._commonService.getAddressTransactions(this.addrHash, this.queryParams.params).subscribe((data: any) => {
      if (data.transactions && data.transactions.length) {
        this.transactions = this.transactions.concat(data.transactions);
        if (data.transactions.length < this.queryParams.limit) {
          this.scrollState = false;
        }
      }
      this._transactionLoading = false;
    });
  }

  onScroll() {
    if (!this._transactionLoading) {
      this.queryParams.next();
    }
  }
}

/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Observable, Subscription} from 'rxjs';
import {filter, tap} from 'rxjs/operators';
import {Params} from '@angular/router/src/shared';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
/*MODELS*/
import {Address} from '../../models/address.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {TOKEN_TYPES} from '../../utils/constants';

@Component({
  selector: 'app-token',
  templateUrl: './token.component.html',
  styleUrls: ['./token.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class TokenComponent implements OnInit, OnDestroy {
  address: Observable<Address>;
  addrHash: string;
  tokenId: string;
  tokenTypes = TOKEN_TYPES;
  apiUrl = this._commonService.getApiUrl();
  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
    private _route: ActivatedRoute,
    private _layoutService: LayoutService,
    ) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._route.params.pipe(
        filter((params: Params) => !!params.id),
      ).subscribe((params: Params) => {
        this.addrHash = params.id;
        this._layoutService.onLoading();
        this.getAddress();
      })
    );
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  getAddress() {
    this.address = this._commonService.getAddress(this.addrHash).pipe(
      filter((addr: Address) => {
        if (!addr || !addr.contract || !addr.go20) {
          this._layoutService.offLoading();
          return false;
        }

        return true;
      }),
      // getting token holder data if address is contract
      tap((addr: Address) => {
        this._layoutService.offLoading();
        addr.ercObj = addr.erc_types.reduce((acc, val) => {
          acc[val] = true;
          return acc;
        }, {});
      })
    );
  }
}

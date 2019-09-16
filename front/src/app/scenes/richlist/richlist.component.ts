/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
import {filter, flatMap, tap} from 'rxjs/operators';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {RichList} from '../../models/rich_list.model';
import {Address} from '../../models/address.model';
import {QueryParams} from '../../models/query_params';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES} from '../../utils/constants';

@Component({
  selector: 'app-richlist',
  templateUrl: './richlist.component.html',
  styleUrls: ['./richlist.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class RichlistComponent implements OnInit, OnDestroy {

  richList: RichList = new RichList();
  richListQueryParams: QueryParams = new QueryParams(50);
  isMoreDisabled = false;
  isLoading = false;

  private _subsArr$: Subscription[] = [];

  static calcSupplyOwned(addresses: Address[], circulatingSupply: any) {
    addresses.forEach((addr: Address) => {
      addr.supplyOwned = (addr.balance / circulatingSupply * 100).toFixed(6);
    });
  }

  constructor(
    private _commonService: CommonService,
    private _layoutService: LayoutService,
    private _metaService: MetaService,
  ) {
    this.initSub();
  }

  ngOnInit(): void {
    this._layoutService.onLoading();
    this.richListQueryParams.init();
    this._metaService.setTitle(META_TITLES.RICHLISLT.title);
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  initSub() {
    this._subsArr$.push(this.richListQueryParams.state.pipe(
      tap(() => this.isLoading = true),
      flatMap(params => this._commonService.getRichlist(params)),
      filter((data: RichList) => !!data),
    ).subscribe((data: RichList) => {
      RichlistComponent.calcSupplyOwned(data.rankings, data.circulating_supply);
      this.richList.rankings = [...this.richList.rankings, ...data.rankings];
      this.richList.circulating_supply = data.circulating_supply;
      this.richList.total_supply = data.total_supply;
      if (data.rankings.length < this.richListQueryParams.limit) {
        this.isMoreDisabled = true;
      }
      this.isLoading = false;
      this._layoutService.offLoading();
    }));
  }
}

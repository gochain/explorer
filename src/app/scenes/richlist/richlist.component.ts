/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
import {flatMap} from 'rxjs/operators';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
/*MODELS*/
import {RichList} from '../../models/rich_list.model';
import {Address} from '../../models/address.model';
import {QueryParams} from '../../models/query_params';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-richlist',
  templateUrl: './richlist.component.html',
  styleUrls: ['./richlist.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class RichlistComponent implements OnInit {

  richList: RichList = new RichList();
  richListQueryParams: QueryParams = new QueryParams(50);
  isMoreDisabled = false;

  private _subsArr$: Subscription[] = [];

  static calcSupplyOwned(addresses: Address[], circulatingSupply: any) {
    addresses.forEach((addr: Address) => {
      addr.supplyOwned = (addr.balance / circulatingSupply * 100).toFixed(6);
    });
  }

  constructor(private _commonService: CommonService, private _layoutService: LayoutService) {
    this.initSub();
  }

  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
    this.richListQueryParams.init();
  }

  initSub() {
    this._subsArr$.push(this.richListQueryParams.state.pipe(
      flatMap(params => this._commonService.getRichlist(params)),
    ).subscribe((data: RichList) => {
      RichlistComponent.calcSupplyOwned(data.rankings, data.circulating_supply);
      this.richList.rankings = [...this.richList.rankings, ...data.rankings];
      this.richList.circulating_supply = data.circulating_supply;
      this.richList.total_supply = data.total_supply;
      if (data.rankings.length < this.richListQueryParams.limit) {
        this.isMoreDisabled = true;
      }
      this._layoutService.isPageLoading.next(false);
    }));
  }
}

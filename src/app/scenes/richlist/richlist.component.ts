/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
/*MODELS*/
import {RichList} from '../../models/rich_list.model';
import {Address} from '../../models/address.model';

@Component({
  selector: 'app-richlist',
  templateUrl: './richlist.component.html',
  styleUrls: ['./richlist.component.css']
})
export class RichlistComponent implements OnInit {

  richList: RichList;

  skip = 0;
  limit = 100;
  isMoreDisabled = false;

  constructor(private _commonService: CommonService, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
    this.richList = new RichList;
    this.richList.rankings = [];
    this.getMore();
  }

  calcSupplyOwned(addresses: Address[], circulatingSupply: any) {
    addresses.forEach((addr: Address) => {
      addr.supplyOwned = (addr.balance / circulatingSupply * 100).toFixed(2);
    });
  }

  getMore() {
    this._commonService.getRichlist(this.skip, this.limit).subscribe((data: RichList) => {
      this.calcSupplyOwned(data.rankings, data.circulating_supply);
      this.richList.rankings = this.richList.rankings.concat(data.rankings);
      this.richList.circulating_supply = data.circulating_supply;
      this.richList.total_supply = data.total_supply;
      this.skip += 100;
      if (data.rankings.length < this.limit) {
        this.isMoreDisabled = true;
      }
      this._layoutService.isPageLoading.next(false);
    });
  }
}

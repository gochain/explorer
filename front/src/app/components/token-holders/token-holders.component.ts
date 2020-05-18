/*CORE*/
import {Component, Input, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Address} from '../../models/address.model';
import {QueryParams} from '../../models/query_params';
import {Holder} from '../../models/holder.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-token-holders',
  templateUrl: './token-holders.component.html',
  styleUrls: ['./token-holders.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class TokenHoldersComponent implements OnInit {
  @Input()
  set addr(value: Address) {
    this._addr = value;
    this.holderQueryParams.setTotalPage(this._addr.number_of_token_holders || 0);
    this.getHolderData();
  }

  get addr(): Address {
    return this._addr;
  }

  token_holders: Holder[] = [];
  holderQueryParams: QueryParams = new QueryParams();

  private _addr: Address;
  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(this.holderQueryParams.state.subscribe(() => {
      this.getHolderData();
    }));
  }

  getHolderData() {
    this._commonService.getAddressHolders(this._addr.address, this.holderQueryParams.params).subscribe((data: any) => {
      this.token_holders = data.token_holders ? data.token_holders.map(x => Object.assign(new Holder(), x)) : [];
    });
  }
}

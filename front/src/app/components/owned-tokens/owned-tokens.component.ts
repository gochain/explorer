/*CORE*/
import {Component, Input, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {QueryParams} from '../../models/query_params';
import {Holder} from '../../models/holder.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-owned-tokens',
  templateUrl: './owned-tokens.component.html',
  styleUrls: ['./owned-tokens.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class OwnedTokensComponent implements OnInit {
  @Input()
  set addrHash(value: string) {
    this._addrHash = value;
    this.getTokenData();
  }

  get addrHash(): string {
    return this._addrHash;
  }

  @Input()
  showPagination = true;


  tokens: Holder[] = [];
  tokensQueryParams: QueryParams = new QueryParams(100);

  private _addrHash: string;
  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(this.tokensQueryParams.state.subscribe(() => {
      this.getTokenData();
    }));
  }

  getTokenData() {
    this._commonService.getAddressTokens(this._addrHash, this.tokensQueryParams.params).subscribe((data: any) => {
      this.tokens = data.owned_tokens ? data.owned_tokens.map(x => Object.assign(new Holder(), x)) : [];
    });
  }
}

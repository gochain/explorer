/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, Params, Router} from '@angular/router';
import {Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {Address} from '../../models/address.model';
import {Contract} from '../../models/contract.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES, TOKEN_TYPES} from '../../utils/constants';

@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class AddressComponent implements OnInit, OnDestroy {
  addr: Address;
  contract: Contract;
  addrHash: string;
  tokenTypes = TOKEN_TYPES;
  apiUrl = this._commonService.getApiUrl();
  tokenId: string;

  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
    private _route: ActivatedRoute,
    private _layoutService: LayoutService,
    private _metaService: MetaService,
    private _router: Router,
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
    this._commonService.getAddress(this.addrHash).pipe(
      filter(value => {
        if (!value) {
          this._layoutService.offLoading();
          return false;
        }
        return true;
      }),
    ).subscribe((addr: Address) => {
      this.addr = addr;
      this._layoutService.offLoading();
      if (this.addr.contract) {
        if (this.addr.token_symbol && this.addr.token_name) {
          this._metaService.setTitle(`${this.addr.token_symbol} - ${this.addr.token_name}`);
        } else {
          this._metaService.setTitle(META_TITLES.CONTRACT.title);
        }
        this.addr.ercObj = this.addr.erc_types.reduce((acc, val) => {
          acc[val] = true;
          return acc;
        }, {});
        this._commonService.getContract(this.addrHash).subscribe(value => {
          this.contract = value;
        });
      } else {
        this._metaService.setTitle(META_TITLES.ADDRESS.title);
      }
    });
  }

  searchToken(): void {
    this._router.navigate([`/token/${this.addrHash}/asset/${this.tokenId}`]);
  }
}

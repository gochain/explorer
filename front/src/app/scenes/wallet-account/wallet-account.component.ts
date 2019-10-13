/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {MetaService} from '../../services/meta.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Address} from '../../models/address.model';
/*UTILS*/
import {META_TITLES} from '../../utils/constants';
import {AutoUnsubscribe} from "../../decorators/auto-unsubscribe";
import {Subscription} from "rxjs";

@Component({
  selector: 'app-wallet-account',
  templateUrl: './wallet-account.component.html',
  styleUrls: ['./wallet-account.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class WalletAccountComponent implements OnInit, OnDestroy {
  addr: Address;
  private _subsArr$: Subscription[] = [];

  constructor(
    public walletService: WalletService,
    private _metaService: MetaService,
    private _commonService: CommonService,
  ) {
  }

  ngOnInit(): void {
    this._metaService.setTitle(META_TITLES.WALLET.title);
    this._subsArr$.push(
      this._commonService.getAddress(this.walletService.account.address).subscribe((addr => {
        this.addr = addr;
      }))
    );
  }

  ngOnDestroy(): void {
    this.walletService.resetProcessing();
  }
}

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

@Component({
  selector: 'app-wallet-account',
  templateUrl: './wallet-account.component.html',
  styleUrls: ['./wallet-account.component.scss']
})
export class WalletAccountComponent implements OnInit, OnDestroy {
  addr: Address;

  constructor(
    public walletService: WalletService,
    private _metaService: MetaService,
    private _commonService: CommonService,
  ) {
  }

  ngOnInit(): void {
    this._metaService.setTitle(META_TITLES.WALLET.title);
    if (!this.walletService.account) {
      return;
    }
    this._commonService.getAddress(this.walletService.account.address).subscribe((addr => {
      this.addr = addr;
    }));
  }

  ngOnDestroy(): void {
    this.walletService.resetProcessing();
  }
}

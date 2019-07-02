/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICES*/
import {WalletService} from '../wallet.service';
import {MetaService} from '../../../services/meta.service';
/*MODELS*/
/*UTILS*/
import {AutoUnsubscribe} from '../../../decorators/auto-unsubscribe';
import {META_TITLES} from '../../../utils/constants';

@Component({
  selector: 'app-wallet-account',
  templateUrl: './wallet-account.component.html',
  styleUrls: ['./wallet-account.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class WalletAccountComponent implements OnInit {
  constructor(
    public walletService: WalletService,
    private _metaService: MetaService,
  ) {
  }

  ngOnInit() {
    this._metaService.setTitle(META_TITLES.WALLET.title);
  }
}

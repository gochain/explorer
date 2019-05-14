/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICES*/
import {WalletService} from '../wallet.service';
import {CommonService} from '../../../services/common.service';
import {MetaService} from '../../../services/meta.service';
/*UTILS*/
import {Account} from 'web3/eth/accounts';
import {META_TITLES} from '../../../utils/constants';

@Component({
  selector: 'app-wallet-create',
  templateUrl: './wallet-create.component.html',
  styleUrls: ['./wallet-create.component.css']
})
export class WalletCreateComponent implements OnInit {

  account: Account;
  apiUrl = this._commonService.getApiUrl();

  constructor(
    private _walletService: WalletService,
    private _commonService: CommonService,
    private _metaService: MetaService,
  ) {
  }

  ngOnInit(): void {
    this._metaService.setTitle(META_TITLES.CREATE_WALLET.title);
    this.account = this._walletService.createAccount();
  }
}

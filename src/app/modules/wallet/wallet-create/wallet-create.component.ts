import {Component, OnInit} from '@angular/core';
import {WalletService} from '../wallet.service';
import {Account} from 'web3/eth/accounts';
import {CommonService} from '../../../services/common.service';

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
  ) {
  }

  ngOnInit(): void {
    this.account = this._walletService.createAccount();
  }
}

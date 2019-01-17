import {Component, OnInit} from '@angular/core';
import {WalletService} from '../wallet.service';
import {Account} from 'web3/eth/accounts';

@Component({
  selector: 'app-wallet-create',
  templateUrl: './wallet-create.component.html',
  styleUrls: ['./wallet-create.component.css']
})
export class WalletCreateComponent implements OnInit {

  account: Account;

  constructor(private _walletService: WalletService) {
  }

  ngOnInit(): void {
    this.account = this._walletService.createAccount();
  }
}

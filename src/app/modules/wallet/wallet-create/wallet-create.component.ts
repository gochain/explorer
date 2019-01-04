import {Component} from '@angular/core';
import {WalletService} from '../wallet.service';

@Component({
  selector: 'app-wallet-create',
  templateUrl: './wallet-create.component.html',
  styleUrls: ['./wallet-create.component.css']
})
export class WalletCreateComponent {

  newAccount: any;

  constructor(private _walletService: WalletService) {
  }

  createAccount(): void {
    this.newAccount = this._walletService.createAccount();
  }
}

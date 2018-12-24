import {Component} from '@angular/core';
import {WalletService} from '../../services/wallet.service';

@Component({
  selector: 'app-create-wallet',
  templateUrl: './create-wallet.component.html',
  styleUrls: ['./create-wallet.component.css']
})
export class CreateWalletComponent {

  newAccount: any;

  constructor(private _walletService: WalletService) {
  }

  createAccount(): void {
    this.newAccount = this._walletService.createAccount();
  }
}

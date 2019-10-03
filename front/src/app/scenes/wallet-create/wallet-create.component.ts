/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {CommonService} from '../../services/common.service';
import {MetaService} from '../../services/meta.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
import {ClipboardService} from 'ngx-clipboard';
/*MODELS*/
import {Account} from 'web3-eth-accounts';
import {PasswordField} from '../../models/password-field.model';
/*UTILS*/
import {META_TITLES} from '../../utils/constants';

@Component({
  selector: 'app-wallet-create',
  templateUrl: './wallet-create.component.html',
  styleUrls: ['./wallet-create.component.css']
})
export class WalletCreateComponent implements OnInit {
  passwordField: PasswordField = new PasswordField();
  account: Account;
  apiUrl = this._commonService.getApiUrl();

  constructor(
    private _walletService: WalletService,
    private _commonService: CommonService,
    private _metaService: MetaService,
    private _router: Router,
    private _toastrService: ToastrService,
    private _clipboardService: ClipboardService,
  ) {
    this._clipboardService.configure({ cleanUpAfterCopy: true });
  }

  ngOnInit(): void {
    this._metaService.setTitle(META_TITLES.CREATE_WALLET.title);
    this._walletService.createAccount().subscribe((account: Account) => {
      this.account = account;
    });
  }

  useWallet(): void {
    this._walletService.openAccount(this.account.privateKey).subscribe((ok: boolean) => {
      if (ok) {
        this._router.navigate(['/wallet/account']);
      }
    });
  }

  onCopy(): void {
    this._toastrService.success('Copied');
  }
}

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
/*UTILS*/
import {META_TITLES} from '../../utils/constants';

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
    private _router: Router,
    private _toastrService: ToastrService,
    private _clipboardService: ClipboardService,
  ) {
    this._clipboardService.configure({ cleanUpAfterCopy: true });
  }

  ngOnInit(): void {
    this._metaService.setTitle(META_TITLES.CREATE_WALLET.title);
    this.account = this._walletService.createAccount();
  }

  useWallet(): void {
    if (this._walletService.openAccount(this.account.privateKey)) {
      this._router.navigate(['/wallet/account']);
    }
  }

  onCopy(): void {
    this._toastrService.success('Copied');
  }
}

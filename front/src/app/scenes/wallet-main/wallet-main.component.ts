/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {FormBuilder, FormControl, FormGroup, Validators} from '@angular/forms';
/*SERVICES*/
import {MetaService} from '../../services/meta.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
import {WalletService} from '../../services/wallet.service';
/*MODELS*/
import {PasswordField} from '../../models/password-field.model';
/*UTILS*/
import {META_TITLES} from '../../utils/constants';
import {LayoutService} from '../../services/layout.service';
import {filter, flatMap} from 'rxjs/operators';

@Component({
  selector: 'app-wallet-main',
  templateUrl: './wallet-main.component.html',
  styleUrls: ['./wallet-main.component.scss']
})
export class WalletMainComponent implements OnInit {
  passwordField: PasswordField = new PasswordField();
  privateKeyForm: FormGroup = this._fb.group({
    privateKey: ['', Validators.compose([Validators.required, WalletMainComponent.checkKeys])],
  });

  static checkKeys(fc: FormControl) {
    if (!fc.value) {
      return;
    }
    const address_or_key = fc.value.toLowerCase();
    if (/^(0x)?[0-9a-f]{40}$/i.test(address_or_key)
      || /^[0-9a-f]{40}$/i.test(address_or_key)
      || /^[0-9a-f]{64}$/i.test(address_or_key)
      || /^(0x)?[0-9a-f]{64}$/i.test(address_or_key)) {
      return null;
    }
    return ({checkKeys: true});
  }

  constructor(
    public walletService: WalletService,
    private _metaService: MetaService,
    private _fb: FormBuilder,
    private _toastrService: ToastrService,
    private _router: Router,
    private _layoutService: LayoutService,
  ) {
  }

  ngOnInit() {
    /*this._layoutService.onLoading();*/
    this._metaService.setTitle(META_TITLES.WALLET.title);
    this.walletService.metamaskConfigured$.pipe(
      filter<boolean>(v => !!v),
      flatMap(() => this.walletService.openAccount()),
    ).subscribe((v) => {
      this._router.navigate(['/wallet/account']);
    });
  }

  onPrivateKeySubmit() {
    const privateKey: string = this.privateKeyForm.get('privateKey').value;
    this.walletService.openAccount(privateKey).subscribe((ok: boolean) => {
      if (ok) {
        this._router.navigate(['/wallet/account']);
      }
    });
  }
}

/*CORE*/
import {Component, OnInit} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
/*UTILS*/
import {DEFAULT_GAS_LIMIT} from '../../utils/constants';

@Component({
  selector: 'app-sender',
  templateUrl: './sender.component.html',
  styleUrls: ['./sender.component.css']
})
export class SenderComponent implements OnInit {

  form: FormGroup = this._fb.group({
    to: ['', Validators.required],
    amount: ['', Validators.required],
    gasLimit: [DEFAULT_GAS_LIMIT, Validators.required],
  });

  constructor(
    private _fb: FormBuilder,
    private _walletService: WalletService,
    private _toastrService: ToastrService,
  ) {
  }

  ngOnInit() {
  }

  sendGo() {
    if (!this.form.valid) {
      this._toastrService.warning('Some field is wrong');
      return;
    }

    const to = this.form.get('to').value;
    const value = this.form.get('amount').value;
    const gas = this.form.get('gasLimit').value;

    this._walletService.sendGo(to, value, gas);
  }
}

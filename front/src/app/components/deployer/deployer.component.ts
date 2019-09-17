/*CORE*/
import {Component, OnInit} from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {Subscription} from 'rxjs';
import {debounceTime, distinctUntilChanged} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
/*MODELS*/
import {TransactionConfig} from 'web3-core';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {isHex} from '../../utils/functions';

@Component({
  selector: 'app-deployer',
  templateUrl: './deployer.component.html',
  styleUrls: ['./deployer.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class DeployerComponent implements OnInit {

  form: FormGroup = this._fb.group({
    byteCode: ['', Validators.required],
    gasLimit: ['', Validators.required],
  });

  private _subsArr$: Subscription[] = [];

  constructor(
    private _fb: FormBuilder,
    private _walletService: WalletService,
    private _toastrService: ToastrService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(this.form.get('byteCode').valueChanges.pipe(
      debounceTime(250),
      distinctUntilChanged(),
    ).subscribe((value: string) => {
      this.estimateDeploymentGas(value);
    }));
  }

  private estimateDeploymentGas(byteCode: string): void {
    if (!byteCode) {
      this.form.get('gasLimit').patchValue('');
      return;
    }
    if (!isHex(byteCode)) {
      this._toastrService.danger('bytecode is not correct');
      return;
    }
    if (!byteCode.startsWith('0x')) {
      byteCode = '0x' + byteCode;
    }
    const tx: TransactionConfig = {data: byteCode};
    this._walletService.estimateGas(tx).pipe(
      // filter((gasLimit: number) => !this.isProcessing),
    ).subscribe((gasLimit: number) => {
      this.form.get('gasLimit').patchValue(gasLimit);
    }, (err) => {
      this._toastrService.danger(err);
      this.form.get('gasLimit').patchValue('');
    });
  }

  deployContract() {
    if (!this.form.valid) {
      this._toastrService.warning('Some field is wrong');
      return;
    }

    const byteCode = this.form.get('byteCode').value;
    const gas = this.form.get('gasLimit').value;

    this._walletService.deployContract(byteCode, gas);
  }
}

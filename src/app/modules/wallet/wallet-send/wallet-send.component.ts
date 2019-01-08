import {Component, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormControl, FormGroup, Validators} from '@angular/forms';
import {WalletService} from '../wallet.service';
import {ToastrService} from '../../toastr/toastr.service';

@Component({
  selector: 'app-wallet-send',
  templateUrl: './wallet-send.component.html',
  styleUrls: ['./wallet-send.component.scss']
})
export class WalletSendComponent implements OnInit {

  privateKeyForm: FormGroup = this._fb.group({
    privateKey: ['', Validators.compose([Validators.required, WalletSendComponent.checkKeys])],
  });

  sendGoForm: FormGroup = this._fb.group({
    to: ['', Validators.required],
    amount: ['', Validators.required],
  });

  deployContractForm: FormGroup = this._fb.group({
    byteCode: ['', Validators.required],
    gasLimit: ['300000', Validators.required],
  });

  useContractForm: FormGroup = this._fb.group({
    contractAddress: ['', Validators.required],
    contractAmount: ['', []],
    contractABI: ['', []],
    contractFunction: [''],
    functionParameters: this._fb.array([]),
  });


  balance: string;
  fromAccount: any;
  address: string; // this is if it's not a private key being used
  receipt: Map<string, any>;
  isSending = false;

  // Contract stuff
  contract: any;
  func: any;
  functionResult: any;
  funcUnsupported: string;

  isOpening = false;

  static checkKeys(fc: FormControl) {
    const address_or_key = fc.value.toLowerCase();
    if (/^(0x)?[0-9a-f]{40}$/i.test(address_or_key)
      || /^[0-9a-f]{40}$/i.test(address_or_key)
      || /^[0-9a-f]{64}$/i.test(address_or_key)
      || /^(0x)?[0-9a-f]{64}$/i.test(address_or_key)) {
      return null;
    }
    return ({checkKeys: true});
  }

  constructor(private _walletService: WalletService, private _fb: FormBuilder, private _toastrService: ToastrService) {
  }

  ngOnInit() {
    this.useContractForm.get('contractAddress').valueChanges.subscribe(val => {
      this.updateContractInfo();
    });
    this.useContractForm.get('contractABI').valueChanges.subscribe(val => {
      this.updateContractInfo();
    });
    this.useContractForm.get('contractFunction').valueChanges.subscribe(val => {
      this.loadFunction();
    });
  }

  loadFunction(): void {
    this.func = null;
    this.functionResult = null;
    this.funcUnsupported = null;
    this.resetFunctionParameter();
    const functionName = this.useContractForm.get('contractFunction').value;
    const abi = this.contract.options.jsonInterface;
    for (let i = 0; i < abi.length; i++) {
      const func = abi[i];
      if (func.name === functionName) {
        this.func = func;
        // TODO: IF ANY INPUTS, add a sub formgroup
        if (func.constant && func.inputs.length === 0) { // if constant, just show value immediately
          // There's a bug in the response here: https://github.com/ethereum/web3.js/issues/1566
          // So doing it myself... :frowning:
          this.callABIFunction(func, []);
        } else {
          // must write a tx to get do this
          if (func.inputs.length > 0) {
            for (const input of func.inputs) {
              this.addFunctionParameter();
            }
            return;
          }

        }
        break;
      }
    }
  }

  get functionParameters() {
    return this.useContractForm.get('functionParameters') as FormArray;
  }

  addFunctionParameter() {
    this.functionParameters.push(this._fb.control(''));
  }

  callABIFunction(func: any, params: string[]): void {
    const m = this.contract.methods[func.name](...params);
    const funcABI = this._walletService.w3.eth.abi.encodeFunctionCall(func, params);
    this._walletService.w3.eth.call({
      to: this.contract.options.address,
      data: '' + funcABI
    }).then(result => {
      const decoded = this._walletService.w3.eth.abi.decodeLog(func.outputs, result, []);
      // This Result object is frikin stupid, it's literaly an empty object that they add fields too
      // convert to something iterable
      const arrR: Array<Array<any>> = new Array<Array<any>>();
      // let mapR: Map<any,any> = new Map<any,any>();
      // for (let j = 0; j < decoded.__length__; j++){
      //   mapR.push([decoded[0], decoded[1]])
      // }
      Object.keys(decoded).forEach(function (key, index) {
        // mapR[key] = decoded[key];
        if (key.startsWith('__')) {
          return;
        }
        arrR.push([key, decoded[key]]);
      });
      this.functionResult = arrR;
    }).catch(err => {
      this._toastrService.danger('ERROR: ' + err);
    });
  }

  resetFunctionParameter() {
    while (this.functionParameters.length !== 0) {
      this.functionParameters.removeAt(0);
    }
  }

  funcsToSelect(): string[] {
    const ret: string[] = [];
    const abi = this.contract.options.jsonInterface;
    // console.log("abi:", abi);
    for (let i = 0; i < abi.length; i++) {
      const func = abi[i];
      if (func.type === 'function') {
        ret.push(func);
      }
    }
    return ret;
  }

  reset() {
    this.balance = null;
    this.fromAccount = null;
    this.address = null;
  }

  closeWallet() {
    this.reset();
    this.privateKeyForm.reset();
  }

  onPrivateKeySubmit() {
    this.reset();
    this.isOpening = true;
    let val: string = this.privateKeyForm.get('privateKey').value;

    if (val.length === 64 && val.indexOf('0x') !== 0) {
      val = '0x' + val;
      this.privateKeyForm.get('privateKey').setValue(val);
    }

    if (val.length === 66) {
      try {
        this.fromAccount = this._walletService.w3.eth.accounts.privateKeyToAccount(val);
        this.address = this.fromAccount.address;
        this.updateBalance();
      } catch (e) {
        this._toastrService.danger('ERROR: ' + e);
        this.isOpening = false;
      }
      return;
    }

    this._toastrService.danger('Given private key is not valid');
    this.isOpening = false;
  }

  updateBalance() {
    if (this._walletService.isAddress(this.address)) {
      this._walletService.getBalance(this.address).subscribe(balance => {
          this._toastrService.info('Updated balance.');
          this.balance = balance;
        },
        err => {
          this._toastrService.danger('ERROR: ' + err);
          this.reset();
        },
        () => this.isOpening = false);
    }
  }

  updateContractInfo(): void {
    const addr: string = this.useContractForm.get('contractAddress').value;
    if (addr.length === 42) {
      // parse the abi
      let abi = this.useContractForm.get('contractABI').value;
      if (abi.length > 0) {
        try {
          abi = JSON.parse(abi);
        } catch (e) {
          return;
        }
        this.contract = new this._walletService.w3.eth.Contract(abi, addr);
        console.log('contract', this.contract);
        console.log('jsonint', this.contract.options.jsonInterface);
      }
    }
  }

  sendGo() {
    if (this.isSending) {
      return;
    }

    this.isSending = true;

    if (!this.sendGoForm.valid) {
      this._toastrService.warning('Some field is wrong');
      this.isSending = false;
      return;
    }

    const to = this.sendGoForm.get('to').value;

    if (to.length !== 42 || !this._walletService.isAddress(to)) {
      this._toastrService.danger('ERROR: Invalid TO address.');
      this.isSending = false;
      return;
    }

    let amount = this.sendGoForm.get('amount').value;

    try {
      amount = this._walletService.w3.utils.toWei(amount, 'ether');
    } catch (e) {
      this._toastrService.danger('ERROR: ' + e);
      this.isSending = false;
      return;
    }

    const tx = {to: to, value: amount, gas: '2000000'};
    const privateKey = this.privateKeyForm.get('privateKey').value;

    this.sendAndWait(privateKey, tx);
  }

  deployContract() {
    if (this.isSending) {
      return;
    }
    this.isSending = true;

    let byteCode = this.deployContractForm.get('byteCode').value;
    if (!byteCode.startsWith('0x')) {
      byteCode = '0x' + byteCode;
    }
    const tx = {data: byteCode, gas: '2000000'};
    const privateKey = this.privateKeyForm.get('privateKey').value;
    this.sendAndWait(privateKey, tx);
  }

  functionName(index) {
    return this.func.inputs[index].name;
  }

  functionPayable(): boolean {
    return this.func && this.func.payable;
  }

  useContract() {
    if (this.isSending) {
      return;
    }
    this.isSending = true;

    const params: string[] = [];
    if (this.func.inputs.length > 0) {
      for (const control of this.functionParameters.controls) {
        params.push(control.value);
      }
    }

    let tx = {};

    const m = this.contract.methods[this.func.name](...params);
    if (this.func.payable) {
      let amount = this.useContractForm.get('contractAmount').value;
      try {
        amount = this._walletService.w3.utils.toWei(amount, 'ether');
      } catch (e) {
        this._toastrService.danger('Cannot convert amount, ERROR: ' + e);
        this.isSending = false;
        return;
      }
      tx = {value: amount};
      Object.assign(tx, tx, {
        to: this.useContractForm.get('contractAddress').value,
        data: m.encodeABI(),
        gas: '2000000'
      });
    } else if (this.func.constant === false) {
      Object.assign(tx, tx, {
        to: this.useContractForm.get('contractAddress').value,
        amount: 0,
        data: m.encodeABI(),
        gas: '2000000'
      });
    } else {
      this.callABIFunction(this.func, params);
      this.isSending = false;
      return;
    }
  }

  sendAndWait(pk: string, tx: any) {
    this._walletService.sendTx(
      pk,
      tx
    ).subscribe(receipt => {
        this.receipt = receipt;
        this.updateBalance();
      },
      err => {
        this._toastrService.danger('ERROR! ' + err);
        this.isSending = false;
      },
      () => {
        this.isSending = false;
        this.resetForms();
      });
  }

  onTabChange() {
    this.receipt = null;
    this.resetForms();
  }

  resetForms() {
    this.sendGoForm.reset();
    this.deployContractForm.reset();
    this.useContractForm.reset();
  }
}

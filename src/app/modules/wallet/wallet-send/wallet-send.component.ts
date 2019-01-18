/*CORE*/
import {Component, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormControl, FormGroup, Validators} from '@angular/forms';
import {Subscription} from 'rxjs';
import {debounceTime, distinctUntilChanged} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../wallet.service';
import {ToastrService} from '../../toastr/toastr.service';
/*MODELS*/
import Contract from 'web3/eth/contract';
import {ABIDefinition} from 'web3/eth/abi';
import {Tx} from 'web3/eth/types';
import {TransactionReceipt} from 'web3/types';
/*UTILS*/
import {AutoUnsubscribe} from '../../../decorators/auto-unsubscribe';

const DEFAULT_GAS_LIMIT = 21000;

@Component({
  selector: 'app-wallet-send',
  templateUrl: './wallet-send.component.html',
  styleUrls: ['./wallet-send.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class WalletSendComponent implements OnInit {

  privateKeyForm: FormGroup = this._fb.group({
    privateKey: ['', Validators.compose([Validators.required, WalletSendComponent.checkKeys])],
  });

  sendGoForm: FormGroup = this._fb.group({
    to: ['', Validators.required],
    amount: ['', Validators.required],
    gasLimit: [DEFAULT_GAS_LIMIT, Validators.required],
  });

  deployContractForm: FormGroup = this._fb.group({
    byteCode: ['', Validators.required],
    gasLimit: [DEFAULT_GAS_LIMIT, Validators.required],
  });

  useContractForm: FormGroup = this._fb.group({
    contractAddress: ['', Validators.required],
    contractAmount: ['', []],
    contractABI: ['', []],
    contractFunction: [''],
    functionParameters: this._fb.array([]),
    gasLimit: [DEFAULT_GAS_LIMIT, Validators.required],
  });


  balance: string;
  fromAccount: any;
  address: string; // this is if it's not a private key being used
  receipt: TransactionReceipt;
  isProcessing = false;

  // Contract stuff
  contract: Contract;
  selectedFunction: ABIDefinition;
  functionResult: any[][];

  isOpening = false;

  private _subsArr$: Subscription[] = [];

  /**
   *
   * @param fc
   */
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

  get functionParameters() {
    return this.useContractForm.get('functionParameters') as FormArray;
  }

  constructor(private _walletService: WalletService, private _fb: FormBuilder, private _toastrService: ToastrService) {
  }

  ngOnInit() {
    this._subsArr$.push(this.useContractForm.get('contractAddress').valueChanges.pipe(
      debounceTime(500),
      distinctUntilChanged(),
    ).subscribe(val => {
      this.updateContractInfo();
    }));
    this._subsArr$.push(this.useContractForm.get('contractABI').valueChanges.pipe(
      debounceTime(500),
      distinctUntilChanged(),
    ).subscribe(val => {
      this.updateContractInfo();
    }));
    this._subsArr$.push(this.useContractForm.get('contractFunction').valueChanges.subscribe(value => {
      this.loadFunction(value);
    }));
  }

  /**
   *
   * @param functionIndex
   */
  loadFunction(functionIndex: number): void {
    this.selectedFunction = null;
    this.functionResult = null;
    this.resetFunctionParameter();
    const abi = this.contract.options.jsonInterface;
    const func = abi[functionIndex];
    this.selectedFunction = func;
    // TODO: IF ANY INPUTS, add a sub formgroup
    // if constant, just show value immediately
    if (func.constant && !func.inputs.length) {
      // There's a bug in the response here: https://github.com/ethereum/web3.js/issues/1566
      // So doing it myself... :frowning:
      this.callABIFunction(func, []);
    } else {
      // must write a tx to get do this
      func.inputs.forEach(() => {
        this.addFunctionParameter();
      });
    }
  }

  addFunctionParameter() {
    this.functionParameters.push(this._fb.control(''));
  }

  /**
   *
   * @param func
   * @param params
   */
  callABIFunction(func: any, params: string[]): void {
    const funcABI: string = this._walletService.w3.eth.abi.encodeFunctionCall(func, params);
    this._walletService.w3.eth.call({
      to: this.contract.options.address,
      data: funcABI,
    }).then((result: string) => {
      const decoded: object = this._walletService.w3.eth.abi.decodeLog(func.outputs, result, []);
      // This Result object is frikin stupid, it's literaly an empty object that they add fields too
      // convert to something iterable
      const arrR: any[][] = [];
      // let mapR: Map<any,any> = new Map<any,any>();
      // for (let j = 0; j < decoded.__length__; j++){
      //   mapR.push([decoded[0], decoded[1]])
      // }
      Object.keys(decoded).forEach((key) => {
        // mapR[key] = decoded[key];
        if (key.startsWith('__')) {
          return;
        }
        arrR.push([key, decoded[key]]);
      });
      this.functionResult = arrR;
    }).catch(err => {
      this._toastrService.danger(err);
    });
  }

  resetFunctionParameter() {
    while (this.functionParameters.length !== 0) {
      this.functionParameters.removeAt(0);
    }
  }

  funcsToSelect(): ABIDefinition[] {
    const abi: ABIDefinition[] = this.contract.options.jsonInterface;
    return abi.filter((abiDef: ABIDefinition) => abiDef.type === 'function');
  }

  reset() {
    this.balance = null;
    this.fromAccount = null;
    this.address = null;
    this.selectedFunction = null;
    this.receipt = null;
  }

  closeWallet() {
    this.reset();
    this.resetForms();
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
        this._toastrService.danger(e);
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
          this.balance = balance.toString();
        },
        err => {
          this._toastrService.danger(err);
          this.isOpening = false;
        },
        () => this.isOpening = false);
    }
  }

  updateContractInfo(): void {
    this.contract = null;
    const addr: string = this.useContractForm.get('contractAddress').value;
    let abi = this.useContractForm.get('contractABI').value;
    if (!addr || !abi) {
      return;
    }
    if (addr.length !== 42) {
      this._toastrService.danger('Wrong contract address');
    }

    if (abi && abi.length > 0) {
      try {
        abi = JSON.parse(abi);
      } catch (e) {
        this._toastrService.danger('Can\'t parse contract abi');
        return;
      }
      try {
        this.contract = new this._walletService.w3.eth.Contract(abi, addr);
      } catch (e) {
        this._toastrService.danger('Can]\'t initiate contract, check entered data');
        return;
      }
    }
  }

  sendGo() {
    if (this.isProcessing) {
      return;
    }

    if (!this.sendGoForm.valid) {
      this._toastrService.warning('Some field is wrong');
      return;
    }

    const to = this.sendGoForm.get('to').value;

    if (to.length !== 42 || !this._walletService.isAddress(to)) {
      this._toastrService.danger('ERROR: Invalid TO address.');
      return;
    }

    let value = this.sendGoForm.get('amount').value;

    try {
      value = this._walletService.w3.utils.toWei(value, 'ether');
    } catch (e) {
      this._toastrService.danger(e);
      return;
    }

    const gas = this.sendGoForm.get('gasLimit').value;

    const tx: Tx = {
      to,
      value,
      gas
    };

    this.sendAndWait(tx);
  }

  deployContract() {
    if (this.isProcessing) {
      return;
    }

    let byteCode = this.deployContractForm.get('byteCode').value;

    if (!byteCode) {
      this._toastrService.danger('ERROR: Invalid data provided.');
      return;
    }

    if (!byteCode.startsWith('0x')) {
      byteCode = '0x' + byteCode;
    }

    const gas = this.deployContractForm.get('gasLimit').value;

    const tx: Tx = {
      data: byteCode,
      gas
    };

    this.sendAndWait(tx);
  }

  useContract() {
    if (this.isProcessing) {
      return;
    }

    const params: string[] = [];

    if (this.selectedFunction.inputs.length) {
      this.functionParameters.controls.forEach(control => {
        params.push(control.value);
      });
    }

    let tx: Tx;

    const m = this.contract.methods[this.selectedFunction.name](...params);
    if (this.selectedFunction.payable) {
      let amount = this.useContractForm.get('contractAmount').value;
      try {
        amount = this._walletService.w3.utils.toWei(amount, 'ether');
      } catch (e) {
        this._toastrService.danger('Cannot convert amount,' + e);
        return;
      }
      tx = {
        to: this.useContractForm.get('contractAddress').value,
        value: amount,
        data: m.encodeABI(),
      };
    } else if (!this.selectedFunction.constant) {
      tx = {
        to: this.useContractForm.get('contractAddress').value,
        data: m.encodeABI(),
      };
    } else {
      this.callABIFunction(this.selectedFunction, params);
      return;
    }

    tx.gas = this.useContractForm.get('gasLimit').value;

    this.sendAndWait(tx);
  }

  sendAndWait(tx: Tx) {
    this.isProcessing = true;

    const privateKey: string = this.privateKeyForm.get('privateKey').value;

    this._walletService.sendTx(
      privateKey,
      tx
    ).subscribe((receipt: TransactionReceipt) => {
        this.receipt = receipt;
        this.updateBalance();
      },
      err => {
        this._toastrService.danger(err);
        this.isProcessing = false;
      });
  }

  onTabChange(tabName: string) {
    /*this.receipt = null;
    switch (tabName) {
      case 'send_go':
        this.sendGoForm.reset();
        break;
      case 'deploy_contract':
        this.deployContractForm.reset();
        break;
      case 'use_contract':
        this.useContractForm.reset();
        break;
    }*/
  }

  resetForms() {
    this.sendGoForm.reset();
    this.sendGoForm.get('gasLimit').setValue(DEFAULT_GAS_LIMIT);
    this.deployContractForm.reset();
    this.deployContractForm.get('gasLimit').setValue(DEFAULT_GAS_LIMIT);
    this.useContractForm.reset();
    this.useContractForm.get('gasLimit').setValue(DEFAULT_GAS_LIMIT);
  }

  resetProcessing() {
    this.isProcessing = false;
    this.receipt = null;
  }
}

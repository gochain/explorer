/*CORE*/
import {Component, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {Subscription} from 'rxjs';
import {debounceTime, distinctUntilChanged} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../wallet.service';
import {ToastrService} from '../../toastr/toastr.service';
/*MODELS*/
import Contract from 'web3/eth/contract';
import {ABIDefinition} from 'web3/eth/abi';

@Component({
  selector: 'app-wallet-use',
  templateUrl: './wallet-use.component.html',
  styleUrls: ['./wallet-use.component.scss']
})
export class WalletUseComponent implements OnInit {

  useContractForm: FormGroup = this._fb.group({
    contractAddress: ['', Validators.required],
    contractAmount: ['', []],
    contractABI: ['', []],
    contractFunction: [''],
    functionParameters: this._fb.array([]),
  });

  // Contract stuff
  contract: Contract;
  selectedFunction: ABIDefinition;
  functionResult: any[][];
  functions: ABIDefinition[] = [];

  isProcessing = false;

  private _subsArr$: Subscription[] = [];

  get functionParameters() {
    return this.useContractForm.get('functionParameters') as FormArray;
  }

  constructor(
    private _walletService: WalletService,
    private _fb: FormBuilder,
    private _toastrService: ToastrService,
  ) {
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
    const func = this.functions[functionIndex];
    this.selectedFunction = func;
    // TODO: IF ANY INPUTS, add a sub formgroup
    // if constant, just show value immediately
    if (func.inputs && func.inputs.length) {
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
    this.isProcessing = true;
    let funcABI: string;
    try {
      funcABI = this._walletService.w3.eth.abi.encodeFunctionCall(func, params);
    } catch (err) {
      this._toastrService.danger(err);
      this.isProcessing = false;
      return;
    }
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
        if (!decoded[key].payable || decoded[key].constant) {
          arrR.push([key, decoded[key]]);
        }
      });
      this.functionResult = arrR;
      this.isProcessing = false;
    }).catch(err => {
      this._toastrService.danger(err);
      this.isProcessing = false;
    });
  }

  resetFunctionParameter() {
    while (this.functionParameters.length !== 0) {
      this.functionParameters.removeAt(0);
    }
  }

  reset() {
    this.selectedFunction = null;
  }

  updateContractInfo(): void {
    this.contract = null;
    this.functions = [];
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
        this.functions = this.contract.options.jsonInterface
          .filter((abiDef: ABIDefinition) => abiDef.type === 'function' && !abiDef.payable && abiDef.constant);
      } catch (e) {
        this._toastrService.danger('Can]\'t initiate contract, check entered data');
        return;
      }
    }
  }

  useContract() {
    const params: string[] = [];

    if (this.selectedFunction.inputs.length) {
      this.functionParameters.controls.forEach(control => {
        params.push(control.value);
      });
    }

    this.callABIFunction(this.selectedFunction, params);
  }
}

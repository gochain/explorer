/*CORE*/
import {Component, Input, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {forkJoin, Observable, Subscription} from 'rxjs';
import {debounceTime, distinctUntilChanged, filter} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Badge} from '../../models/badge.model';
import {Address} from '../../models/address.model';
import {Contract} from '../../models/contract.model';
import Web3 from 'web3';
import {TransactionConfig} from 'web3-core';
import {AbiItem, toWei} from 'web3-utils';
import {Contract as Web3Contract} from 'web3-eth-contract';
/*UTILS*/
import {ErcName} from '../../utils/enums';
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {ContractAbi} from '../../utils/types';
import {getAbiMethods, getDecodedData, makeContractAbi, makeContractBadges} from '../../utils/functions';
import {ERC_INTERFACE_IDENTIFIERS} from '../../utils/constants';
import BigNumber from 'bignumber.js';

@Component({
  selector: 'app-interactor',
  templateUrl: './interactor.component.html',
  styleUrls: ['./interactor.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class InteractorComponent implements OnInit {

  form: FormGroup = this._fb.group({
    contractAddress: ['', Validators.required],
    contractAmount: ['', []],
    contractABI: ['', Validators.required],
    contractFunction: [''],
    functionParameters: this._fb.array([]),
    gasLimit: [''],
    erc: [''],
  });

  contractBadges: Badge[] = [];
  abiTemplates = [ErcName.Go20, ErcName.Go721, ErcName.Go165, ErcName.Upgradeable, ErcName.AllFunctions];

  contract: Web3Contract;
  abiFunctions: AbiItem[];
  selectedFunction: AbiItem;
  functionResult: { output: any[][], error: string };
  addr: Address;

  hasData = false;

  @Input('contractData')
  set address([addr, contract]: [Address, Contract]) {
    this.hasData = true;
    this.form.patchValue({
      contractAddress: addr.address,
    }, {
      emitEvent: false,
    });
    if (contract) {
      this.handleContractData(addr, contract);
    }
  }

  private _subsArr$: Subscription[] = [];

  get functionParameters() {
    return this.form.get('functionParameters') as FormArray;
  }

  constructor(
    public walletService: WalletService,
    private _fb: FormBuilder,
    private _toastrService: ToastrService,
    private _commonService: CommonService,
    private _activatedRoute: ActivatedRoute,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._activatedRoute.queryParamMap.pipe(
        filter((params: ParamMap) => params.has('address'))
      ).subscribe((params: ParamMap) => {
        const addr = params.get('address');
        if (addr.length === 42) {
          this.form.patchValue({
            contractAddress: addr
          }, {
            emitEvent: true,
          });
        } else {
          this._toastrService.warning('Contract address is invalid');
        }
      })
    );
    this._subsArr$.push(this.form.get('contractAddress').valueChanges.pipe(
      debounceTime(500),
      distinctUntilChanged(),
    ).subscribe((val: string) => {
      this.updateContract();
      this.getContractData(val);
    }));
    this._subsArr$.push(this.form.get('contractABI').valueChanges.pipe(
      debounceTime(500),
      distinctUntilChanged(),
    ).subscribe(val => {
      this.updateContract();
    }));
    this._subsArr$.push(this.form.get('contractFunction').valueChanges.subscribe((value: number) => {
      this.onDefinitionSelect(value);
    }));
    this._subsArr$.push((this.form.get('functionParameters') as FormArray).valueChanges.pipe(
      debounceTime(1200),
      distinctUntilChanged(),
    ).subscribe((values: string[]) => {
      this.estimateFunctionGas(values);
    }));
    this._subsArr$.push(this.form.get('erc').valueChanges.subscribe(value => {
      this.onAbiTemplateSelect(value);
    }));
  }

  /**
   *
   * @param functionIndex
   */
  onDefinitionSelect(functionIndex: number): void {
    this.selectedFunction = null;
    this.functionResult = null;
    this.functionParameters.clear();
    this.selectedFunction = this.abiFunctions[functionIndex];
    // TODO: IF ANY INPUTS, add a sub formgroup
    // if constant, just show value immediately
    if (this.selectedFunction.constant && !this.selectedFunction.inputs.length) {
      // There's a bug in the response here: https://github.com/ethereum/web3.js/issues/1566
      // So doing it myself... :frowning:
      this.callABIFunction(this.selectedFunction);
    } else {
      // must write a tx to get do this
      this.selectedFunction.inputs.forEach(() => {
        this.functionParameters.push(this._fb.control(''));
      });
    }
  }


  onTokenValueChange(event, controlIndex: number): void {
    let value: string = (<HTMLInputElement>event.target).value;
    if (value) {
      value = (new BigNumber(value)).multipliedBy('1e' + this.addr.decimals).toString();
      if (/e+/.test(value)) {
        const parts = value.split('e+');
        let first = parts[0].replace('.', '');
        const zeroes = parseInt(parts[1], 10) - (first.length - 1);
        for (let i = 0; i < zeroes; i++) {
          first += '0';
        }
        value = first;
      }
    }
    this.functionParameters.controls[controlIndex].patchValue(value, {
      emitEvent: true,
    });
  }

  /**
   *
   * @param func
   * @param params
   */
  callABIFunction(func: AbiItem, params: string[] = []): void {
    this.walletService.call(this.contract.options.address, func, params).subscribe((decoded: object) => {
      if (!decoded) {
        this.functionResult = {error: 'Result is empty', output: null};
        return;
      }
      this.functionResult = {output: getDecodedData(decoded, func, this.addr), error: null};
    }, err => {
      this.functionResult = {error: err, output: null};
    });
  }

  private getContractData(addrHash: string) {
    if (!addrHash && addrHash.length !== 42) {
      return;
    }
    forkJoin<Address, Contract>([
      this._commonService.getAddress(addrHash),
      this._commonService.getContract(addrHash),
    ]).pipe(
      filter((data: [Address, Contract]) => !!data[0] && !!data[1]),
    ).subscribe((data: [Address, Contract]) => {
      this.handleContractData(data[0], data[1]);
    });
  }

  private handleContractData(address: Address, contract: Contract) {
    this.addr = address;
    this.contractBadges = makeContractBadges(address, contract);
    if (contract.abi && contract.abi.length) {
      this.form.patchValue({
        contractABI: JSON.stringify(contract.abi, null, 2),
      }, {
        emitEvent: true,
      });
    } else if (address.interfaces && address.interfaces.length) {
      this._commonService.abi$.subscribe((abiDefinitions: ContractAbi) => {
        const abi: AbiItem[] = address.interfaces.reduce((acc, abiName) => {
          if (abiDefinitions[abiName]) {
            acc.push(abiDefinitions[abiName]);
          }
          return acc;
        }, []);
        this.form.patchValue({
          contractABI: JSON.stringify(abi, null, 2),
        }, {
          emitEvent: true,
        });
      });
    }
  }

  private estimateFunctionGas(values: string[]): void {
    if (!this.selectedFunction.payable && this.selectedFunction.constant) {
      return;
    }
    if (values.some(value => !value)) {
      this.form.get('gasLimit').patchValue('');
      return;
    }
    let tx: TransactionConfig;

    try {
      tx = this.formTx(values);
    } catch (e) {
      return;
    }

    this.walletService.estimateGas(tx).pipe(
      // filter((gasLimit: number) => !this.isProcessing),
    ).subscribe((gasLimit: number) => {
      this.form.get('gasLimit').patchValue(gasLimit);
    }, (err) => {
      this._toastrService.danger(err);
      this.form.get('gasLimit').patchValue('');
    });
  }

  formTx(params: string[]): TransactionConfig {
    const m = this.contract.methods[this.selectedFunction.name](...params);

    const tx: TransactionConfig = {
      to: this.contract.options.address,
      data: m.encodeABI(),
      from: this.walletService.accountAddress,
    };

    if (this.selectedFunction.payable) {
      const amount = this.form.get('contractAmount').value;
      try {
        tx.value = toWei(amount, 'ether').toString(10);
      } catch (e) {
        throw Error('Cannot convert amount,' + e);
      }
    }
    return tx;
  }

  updateContract(): void {
    const addrHash: string = this.form.get('contractAddress').value;
    const abi: string = this.form.get('contractABI').value;
    if (!addrHash || !abi) {
      return;
    }
    if (addrHash.length !== 42) {
      this._toastrService.danger('Wrong contract address');
      return;
    }

    this.contract = null;
    let abiItem: AbiItem[];

    try {
      abiItem = JSON.parse(abi);
    } catch (e) {
      this._toastrService.danger('Can\'t parse contract abi');
      return;
    }

    this._initContract(addrHash, abiItem);
  }

  useContract(): void {
    const params: string[] = [];

    if (this.selectedFunction.inputs.length) {
      this.functionParameters.controls.forEach(control => {
        params.push(control.value);
      });
    }

    let tx: TransactionConfig;

    if (this.selectedFunction.payable || !this.selectedFunction.constant) {
      try {
        tx = this.formTx(params);
      } catch (e) {
        this._toastrService.danger(e);
        return;
      }
    } else {
      this.callABIFunction(this.selectedFunction, params);
      return;
    }

    tx.gas = this.form.get('gasLimit').value;
    this.walletService.sendTx(tx);
  }

  onAbiTemplateSelect(ercName: ErcName) {
    this._commonService.abi$.subscribe((abi: ContractAbi) => {
      const ABI: AbiItem[] = makeContractAbi(ERC_INTERFACE_IDENTIFIERS[ercName], abi);
      this.form.patchValue({
        contractABI: JSON.stringify(ABI, null, 2),
      }, {
        emitEvent: true,
      });
    });
  }

  private _initContract(addrHash: string, abiItems: AbiItem[]) {
    this.walletService.initContract(addrHash, abiItems).subscribe(
      (contract) => {
        this.contract = contract;
        this.abiFunctions = getAbiMethods(abiItems);
      }, (err) => {
        this._toastrService.danger('Can\'t initiate contract, check entered data');
        console.error(`Failed to initiate contract (${addrHash}): ${err}`);
      }
    );
  }
}

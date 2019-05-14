/*CORE*/
import {Component, Input, OnInit} from '@angular/core';
import {FormArray, FormBuilder, FormGroup, Validators} from '@angular/forms';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {forkJoin, Subscription} from 'rxjs';
import {debounceTime, distinctUntilChanged, filter} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../wallet.service';
import {ToastrService} from '../../toastr/toastr.service';
import {CommonService} from '../../../services/common.service';
import {MetaService} from '../../../services/meta.service';
/*MODELS*/
import Web3Contract from 'web3/eth/contract';
import {ABIDefinition} from 'web3/eth/abi';
import {Contract} from '../../../models/contract.model';
import {Badge} from '../../../models/badge.model';
import {Address} from '../../../models/address.model';
/*UTILS*/
import {ErcName} from '../../../utils/enums';
import {ERC_INTERFACE_IDENTIFIERS, META_TITLES} from '../../../utils/constants';
import {getAbiMethods, getDecodedData, makeContractAbi, makeContractBadges} from '../../../utils/functions';
import {ContractAbi} from '../../../utils/types';

@Component({
  selector: 'app-wallet-use',
  templateUrl: './wallet-use.component.html',
  styleUrls: ['./wallet-use.component.scss'],
})
export class WalletUseComponent implements OnInit {

  hasData = false;

  @Input('contractData')
  set address([addr, contract]: [Address, Contract]) {
    this.hasData = true;
    this.useContractForm.patchValue({
      contractAddress: addr.address,
    }, {
      emitEvent: false,
    });
    if (contract) {
      this.handleContractData(addr, contract);
    }
  }

  useContractForm: FormGroup = this._fb.group({
    contractAddress: ['', Validators.required],
    contractAmount: ['', []],
    contractABI: ['', []],
    contractFunction: [''],
    functionParameters: this._fb.array([]),
  });

  // Contract stuff
  contract: Web3Contract;
  selectedFunction: ABIDefinition;
  functionResult: any[][];
  functions: ABIDefinition[] = [];

  isProcessing = false;

  contractBadges: Badge[] = [];

  abiTemplates = [ErcName.Go20, ErcName.Go721];

  addr: Address;

  private _subsArr$: Subscription[] = [];

  get functionParameters() {
    return this.useContractForm.get('functionParameters') as FormArray;
  }

  constructor(
    private _walletService: WalletService,
    private _fb: FormBuilder,
    private _toastrService: ToastrService,
    private _activatedRoute: ActivatedRoute,
    private _commonService: CommonService,
    private metaService: MetaService,
  ) {
  }

  ngOnInit() {
    this.metaService.setTitle(META_TITLES.USE_CONTRACT.title);
    this._subsArr$.push(
      this._activatedRoute.queryParamMap.pipe(
        filter((params: ParamMap) => params.has('address'))
      ).subscribe((params: ParamMap) => {
        const addr = params.get('address');
        if (addr.length === 42) {
          this.useContractForm.patchValue({
            contractAddress: addr
          });
          this.getContractData(addr);
        } else {
          this._toastrService.warning('Contract address is invalid');
        }
      })
    );
    this._subsArr$.push(this.useContractForm.get('contractAddress').valueChanges.pipe(
      debounceTime(500),
      distinctUntilChanged(),
    ).subscribe((val: string) => {
      this.updateContractInfo();
      this.getContractData(val);
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

  private getContractData(addrHash: string) {
    forkJoin([
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
      this.useContractForm.patchValue({
        contractABI: JSON.stringify(contract.abi),
      }, {
        emitEvent: false,
      });
      this.initiateContract(contract.abi, address.address);
    } else if (address.interfaces && address.interfaces.length) {
      this._walletService.abi$.subscribe((abiDefinitions: ContractAbi) => {
        const abi: ABIDefinition[] = address.interfaces.reduce((acc, abiName) => {
          if (abiDefinitions[abiName]) {
            acc.push(abiDefinitions[abiName]);
          }
          return acc;
        }, []);
        this.useContractForm.patchValue({
          contractABI: JSON.stringify(abi),
        }, {
          emitEvent: false,
        });
        this.initiateContract(abi, address.address);
      });
    }
  }

  private initiateContract(abi: ABIDefinition[], addrHash: string) {
    try {
      this.contract = new this._walletService.w3.eth.Contract(abi, addrHash);
      this.functions = getAbiMethods(this.contract.options.jsonInterface);
    } catch (e) {
      this._toastrService.danger('Can]\'t initiate contract, check entered data');
      return;
    }
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
    if (func && func.inputs && func.inputs.length) {
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
  callABIFunction(func: ABIDefinition, params: string[]): void {
    this.isProcessing = true;
    this._walletService.call(this.contract.options.address, func, params).then((decoded: object) => {
      this.functionResult = getDecodedData(decoded, func, this.addr);
    }).catch(err => {
      this._toastrService.danger(err);
    }).then(() => {
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
    if (!addr) {
      return;
    }
    if (addr.length !== 42) {
      this._toastrService.danger('Wrong contract address');
      return;
    }
    if (!abi) {
      return;
    }

    if (abi && abi.length > 0) {
      try {
        abi = JSON.parse(abi);
      } catch (e) {
        this._toastrService.danger('Can\'t parse contract abi');
        return;
      }
      this.initiateContract(abi, addr);
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

  onAbiTemplateClick(ercName: ErcName) {
    this._walletService.abi$.subscribe((abi: ContractAbi) => {
      const ABI: ABIDefinition[] = makeContractAbi(ERC_INTERFACE_IDENTIFIERS[ercName], abi);
      const addr: string = this.useContractForm.get('contractAddress').value;
      this.useContractForm.patchValue({
        contractABI: JSON.stringify(ABI),
      }, {
        emitEvent: false,
      });
      if (addr.length === 42 && ABI.length) {
        this.initiateContract(ABI, addr);
      }
    });
  }
}

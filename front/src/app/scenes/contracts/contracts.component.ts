/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {FormBuilder, FormGroup} from '@angular/forms';
import {Subscription} from 'rxjs';
import {flatMap, tap} from 'rxjs/operators';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {MetaService} from '../../services/meta.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
/*MODELS*/
import {Address} from '../../models/address.model';
import {QueryParams} from '../../models/query_params';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES, TOKEN_TYPES} from '../../utils/constants';

interface ISelectOption {
  value: string | number;
  desc: string | number;
}

const TOKEN_TYPE_OPTIONS: ISelectOption[] = Object.keys(TOKEN_TYPES).map((key) => ({value: key, desc: TOKEN_TYPES[key]}));

const SORT_FIELD_OPTIONS: ISelectOption[] = [
  {
    value: 'number_of_transactions',
    desc: 'Number of transactions',
  }, {
    value: 'number_of_token_holders',
    desc: 'Number of token holders',
  }, {
    value: 'number_of_internal_transactions',
    desc: 'Number of internal transactions'
  }, {
    value: 'number_of_token_transactions',
    desc: 'Number of token transactions'
  },
];

@Component({
  selector: 'app-contracts',
  templateUrl: './contracts.component.html',
  styleUrls: ['./contracts.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class ContractsComponent implements OnInit, OnDestroy {

  tokenTypes = TOKEN_TYPES;
  tokenTypeOptions = TOKEN_TYPE_OPTIONS;
  sortFieldOptions = SORT_FIELD_OPTIONS;
  addresses: Address[] = [];
  contractsQueryParams: QueryParams = new QueryParams(50);
  isMoreDisabled = false;
  isLoading = false;

  filter: FormGroup = this._fb.group({
    contract_name: [''],
    token_name: [''],
    token_symbol: [''],
    erc_type: [''],
    sortby: [''],
    asc: [false],
  });

  private _subsArr$: Subscription[] = [];

  constructor(
    private _fb: FormBuilder,
    private _commonService: CommonService,
    private _layoutService: LayoutService,
    private _metaService: MetaService,
    private _toastrService: ToastrService,
  ) {
    this.initSub();
  }

  ngOnInit() {
    this._layoutService.onLoading();
    this.contractsQueryParams.init();
    this._metaService.setTitle(META_TITLES.CONTRACTS.title);
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  initSub() {
    this._subsArr$.push(this.contractsQueryParams.state.pipe(
      tap(() => this.isLoading = true),
      flatMap(params => this._commonService.getContractsList(params)),
    ).subscribe((data: Address[]) => {
      this.addresses = this.contractsQueryParams.skip === 0 ? data : [...this.addresses, ...data];
      if (data.length < this.contractsQueryParams.limit) {
        this.isMoreDisabled = true;
      }
      this.isLoading = false;
      this._layoutService.offLoading();
    }));
  }

  onFilterSubmit(): void {
    this.contractsQueryParams.filter = this.filter.getRawValue();
  }
}

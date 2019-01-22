/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Observable, Subscription} from 'rxjs';
import {filter, tap} from 'rxjs/operators';
import {Params} from '@angular/router/src/shared';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
/*MODELS*/
import {Address} from '../../models/address.model';
import {Transaction} from '../../models/transaction.model';
import {Holder} from '../../models/holder.model';
import {QueryParams} from '../../models/query_params';
import {InternalTransaction} from '../../models/internal-transaction.model';
import {Contract} from '../../models/contract.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';


@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class AddressComponent implements OnInit, OnDestroy {
  address: Observable<Address>;
  transactions: Transaction[] = [];
  token_holders: Holder[] = [];
  internal_transactions: InternalTransaction[] = [];
  contract: Contract;
  transactionQueryParams: QueryParams = new QueryParams();
  internalTransactionQueryParams: QueryParams = new QueryParams();
  holderQueryParams: QueryParams = new QueryParams();
  addrHash: string;
  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService, private _route: ActivatedRoute, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._route.params.pipe(
        filter((params: Params) => !!params.id),
      ).subscribe((params: Params) => {
        this.transactions = [];
        this.addrHash = params.id;
        this._layoutService.onLoading();
        this.getAddress();
      })
    );
    this._subsArr$.push(this.transactionQueryParams.state.subscribe(() => {
      this.getTransactionData();
    }));
    this._subsArr$.push(this.holderQueryParams.state.subscribe(() => {
      this.getHolderData();
    }));
    this._subsArr$.push(this.internalTransactionQueryParams.state.subscribe(() => {
      this.getInternalTransactions();
    }));
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  getAddress() {
    this.address = this._commonService.getAddress(this.addrHash).pipe(
      filter(value => {
        if (!value) {
          this._layoutService.offLoading();
          return false;
        }

        return true;
      }),
      // getting token holder data if address is contract
      tap((addr: Address) => {
        this._layoutService.offLoading();
        this.transactionQueryParams.setTotalPage(addr.number_of_transactions);
        if (addr.contract && addr.go20) {
          this.holderQueryParams.setTotalPage(addr.number_of_token_holders);
          this.internalTransactionQueryParams.setTotalPage(addr.number_of_internal_transactions);
          this.getHolderData();
          this.getInternalTransactions();
        }
        if (addr.contract) {
          this.getContractData();
        }
        this.getTransactionData();
      })
    );
  }

  getTransactionData() {
    this._commonService.getAddressTransactions(this.addrHash, this.transactionQueryParams.params).subscribe((data: any) => {
      this.transactions = data.transactions;
    });
  }

  getHolderData() {
    this._commonService.getAddressHolders(this.addrHash, this.holderQueryParams.params).subscribe((data: any) => {
      this.token_holders = data.token_holders;
    });
  }

  getInternalTransactions() {
    this._commonService.getAddressInternalTransaction(this.addrHash, this.internalTransactionQueryParams.params).subscribe((data: any) => {
      this.internal_transactions = data.internal_transactions;
    });
  }

  getContractData() {
    this._commonService.getContract(this.addrHash).subscribe((data: Contract) => {
      this.contract = data;
    });
  }
}

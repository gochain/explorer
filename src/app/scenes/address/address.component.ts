/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {forkJoin, Observable, of, Subscription} from 'rxjs';
import {concatMap, filter, map} from 'rxjs/operators';
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
import {TOKEN_TYPES} from '../../utils/constants';

@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class AddressComponent implements OnInit, OnDestroy {
  addr: Address;
  transactions: Transaction[] = [];
  token_holders: Holder[] = [];
  // address owned tokens
  tokens: Holder[] = [];
  internal_transactions: InternalTransaction[] = [];
  token_transactions: InternalTransaction[] = [];
  contract: Contract;
  transactionQueryParams: QueryParams = new QueryParams();
  internalTransactionQueryParams: QueryParams = new QueryParams();
  tokenTransactionQueryParams: QueryParams = new QueryParams();
  holderQueryParams: QueryParams = new QueryParams();
  tokensQueryParams: QueryParams = new QueryParams(100);
  addrHash: string;
  tokenTypes = TOKEN_TYPES;
  apiUrl = this._commonService.getApiUrl();
  tokenId: string;

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
    this._subsArr$.push(this.tokensQueryParams.state.subscribe(() => {
      this.getTokenData();
    }));
    this._subsArr$.push(this.internalTransactionQueryParams.state.subscribe(() => {
      this.getInternalTransactions();
    }));
    this._subsArr$.push(this.tokenTransactionQueryParams.state.subscribe(() => {
      this.getTokenTransactions();
    }));
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  getAddress() {
    this._commonService.getAddress(this.addrHash).pipe(
      filter(value => {
        if (!value) {
          this._layoutService.offLoading();
          return false;
        }

        return true;
      }),
    ).subscribe((addr: Address) => {
      this.addr = addr;
      this._layoutService.offLoading();
      this.transactionQueryParams.setTotalPage(addr.number_of_transactions || 0);
      this.getTransactionData();
      this.tokenTransactionQueryParams.setTotalPage(addr.number_of_token_transactions || 0);
      this.getTokenTransactions();
      if (this.addr.contract) {
        this.holderQueryParams.setTotalPage(addr.number_of_token_holders || 0);
        this.internalTransactionQueryParams.setTotalPage(addr.number_of_internal_transactions || 0);
        this.getHolderData();
        this.getInternalTransactions();
        this.addr.ercObj = this.addr.erc_types.reduce((acc, val) => {
          acc[val] = true;
          return acc;
        }, {});

        this.getContractData();
      } else {
        this.tokensQueryParams.setTotalPage(100);
        this.getTokenData();
      }
    });
  }

  getTransactionData() {
    this._commonService.getAddressTransactions(this.addrHash, this.transactionQueryParams.params).subscribe((data: any) => {
      this.transactions = data.transactions || [];
    });
  }

  getHolderData() {
    this._commonService.getAddressHolders(this.addrHash, this.holderQueryParams.params).subscribe((data: any) => {
      this.token_holders = data.token_holders || [];
    });
  }

  getTokenData() {
    this._commonService.getAddressTokens(this.addrHash, this.tokensQueryParams.params).subscribe((data: any) => {
      this.tokens = data.owned_tokens || [];
    });
  }

  getInternalTransactions() {
    this._commonService.getAddressInternalTransaction(this.addrHash, {
      ...this.internalTransactionQueryParams.params,
      token_transactions: false,
    }).subscribe((data: any) => {
      this.internal_transactions = data.internal_transactions || [];
    });
  }

  getTokenTransactions() {
    this._commonService.getAddressInternalTransaction(this.addrHash, {
      ...this.tokenTransactionQueryParams.params,
      token_transactions: true,
    }).pipe(
      concatMap((data: any) => {
        if (!data.internal_transactions || !data.internal_transactions.length) {
          return of(null);
        }
        const contractAddresses: string[] = [];
        data.internal_transactions.forEach((tx: InternalTransaction) => {
          if (!this._commonService.contractsCache[tx.contract_address]) {
            contractAddresses.push(tx.contract_address);
          }
        });
        return forkJoin<Observable<any>[]>(contractAddresses.map((addr: string) => {
          return this._commonService.getAddress(addr);
        })).pipe(
          map((addrs: Address[]) => {
            addrs.forEach((item: Address) => {
              this._commonService.contractsCache[item.address] = item;
            });
            data.internal_transactions.forEach((tx: InternalTransaction) => {
              tx.address = this._commonService.contractsCache[tx.contract_address];
            });
            return data.internal_transactions;
          })
        );
      })
    ).subscribe((data: any) => {
      this.token_transactions = data || [];
    });
  }

  getContractData() {
    this._commonService.getContract(this.addrHash).subscribe((data: Contract) => {
      this.contract = data;
    });
  }
}

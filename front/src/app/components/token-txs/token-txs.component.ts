/*CORE*/
import { Component, Input, OnInit } from '@angular/core';
import { forkJoin, of, Subscription } from 'rxjs';
import { concatMap, map } from 'rxjs/operators';
/*SERVICES*/
import { CommonService } from '../../services/common.service';
/*MODELS*/
import { InternalTransaction } from '../../models/internal-transaction.model';
import { QueryParams } from '../../models/query_params';
import { Address } from '../../models/address.model';
/*UTILS*/
import { AutoUnsubscribe } from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-token-txs',
  templateUrl: './token-txs.component.html',
  styleUrls: ['./token-txs.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class TokenTxsComponent implements OnInit {
  @Input()
  set addr(value: Address) {
    this._addr = value;
    this.tokenTransactionQueryParams.setTotalPage(this._addr.number_of_token_transactions || 0);
    this.getTokenTransactions();
  }

  get addr(): Address {
    return this._addr;
  }

  token_transactions: InternalTransaction[] = [];
  tokenTransactionQueryParams: QueryParams = new QueryParams();

  private _addr: Address;
  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(this.tokenTransactionQueryParams.state.subscribe(() => {
      this.getTokenTransactions();
    }));
  }

  getTokenTransactions() {
    this._commonService.getAddressInternalTransaction(this.addr.address, {
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
        if (contractAddresses.length) {
          return forkJoin<any>(contractAddresses.map((addr: string) => {
            return this._commonService.getAddress(addr);
          })).pipe(
            map((addrs: Address[]) => {
              addrs.forEach((item: Address) => {
                this._commonService.contractsCache[item.address] = item;
              });
              return data.internal_transactions;
            })
          );
        } else {
          return of(data.internal_transactions);
        }
      }),
      map((data: InternalTransaction[]) => {
        if (data) {
          data.forEach((tx: InternalTransaction) => {
            tx.address = this._commonService.contractsCache[tx.contract_address];
          });
        }
        return data;
      }),
    ).subscribe((data: InternalTransaction[]) => {
      this.token_transactions = data || [];
    });
  }
}

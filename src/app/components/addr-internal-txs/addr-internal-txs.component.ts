/*CORE*/
import {Component, Input, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {InternalTransaction} from '../../models/internal-transaction.model';
import {QueryParams} from '../../models/query_params';
import {Address} from '../../models/address.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-addr-internal-txs',
  templateUrl: './addr-internal-txs.component.html',
  styleUrls: ['./addr-internal-txs.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class AddrInternalTxsComponent implements OnInit {
  @Input()
  set addr(value: Address) {
    this._addr = value;
    this.internalTransactionQueryParams.setTotalPage(this._addr.number_of_transactions || 0);
    this.getInternalTransactions();
  }

  get addr(): Address {
    return this._addr;
  }

  internal_transactions: InternalTransaction[] = [];
  internalTransactionQueryParams: QueryParams = new QueryParams();

  private _addr: Address;
  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(this.internalTransactionQueryParams.state.subscribe(() => {
      this.getInternalTransactions();
    }));
  }

  getInternalTransactions() {
    this._commonService.getAddressInternalTransaction(this.addr.address, {
      ...this.internalTransactionQueryParams.params,
      token_transactions: false,
    }).subscribe((data: any) => {
      this.internal_transactions = data.internal_transactions || [];
    });
  }
}

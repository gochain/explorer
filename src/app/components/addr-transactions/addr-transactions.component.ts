import {Component, Input, OnInit} from '@angular/core';
import {Transaction} from '../../models/transaction.model';
import {CommonService} from '../../services/common.service';
import {Subscription} from 'rxjs';
import {QueryParams} from '../../models/query_params';
import {Address} from '../../models/address.model';

@Component({
  selector: 'app-addr-transactions',
  templateUrl: './addr-transactions.component.html',
  styleUrls: ['./addr-transactions.component.css']
})
export class AddrTransactionsComponent implements OnInit {
  @Input()
  set addr(value: Address) {
    this._addr = value;
    this.transactionQueryParams.setTotalPage(this._addr.number_of_transactions || 0);
    this.getTransactionData();
  }

  get addr(): Address {
    return this._addr;
  }

  transactions: Transaction[] = [];
  transactionQueryParams: QueryParams = new QueryParams();

  private _addr: Address;
  private _subsArr$: Subscription[] = [];

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(this.transactionQueryParams.state.subscribe(() => {
      this.getTransactionData();
    }));
  }

  getTransactionData() {
    this._commonService.getAddressTransactions(this._addr.address, this.transactionQueryParams.params).subscribe((data: any) => {
      this.transactions = data.transactions || [];
    });
  }
}

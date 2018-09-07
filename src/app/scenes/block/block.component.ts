/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
import {ActivatedRoute} from '@angular/router';
import {Params} from '@angular/router/src/shared';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/template.service';
/*MODELS*/
import {Block} from '../../models/block.model';
import {QueryParams} from '../../models/query_params';
import {Transaction} from '../../models/transaction.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class BlockComponent implements OnInit {
  block: Block;
  transactions: Transaction[] = [];
  transactionQueryParams: QueryParams = new QueryParams();

  private _blockNum: number;
  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService, private _route: ActivatedRoute, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._subsArr$.push(this._route.params.pipe(
      filter((params: Params) => !!params.id),
    ).subscribe((params: Params) => {
      this._blockNum = params.id;
      this._layoutService.isPageLoading.next(true);
      this.getData();
    }));
    this._subsArr$.push(this.transactionQueryParams.state.subscribe(() => {
      this.getData();
    }));
  }

  getData() {
    this._commonService.getBlock(this._blockNum, this.transactionQueryParams.params).subscribe((data: Block) => {
      this.block = data;
      this.transactionQueryParams.setTotalPage(this.block.tx_count);
      if (this.block.tx_count) {
        this.getTransactionData();
      } else {
        this.transactions = [];
      }
      this._layoutService.isPageLoading.next(false);
    });
  }

  // to-do: add caching
  getTransactionData() {
    this._commonService.getBlockTransactions(this._blockNum, this.transactionQueryParams.params).subscribe((data: any) => {
      this.transactions = data.transactions;
    });
  }

  onTransactionPageSelect(page: number) {
    this.transactionQueryParams.toPage(page);
  }
}

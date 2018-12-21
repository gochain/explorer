/*CORE*/
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Observable } from 'rxjs';
import { switchMap, tap } from 'rxjs/operators';
/*SERVICES*/
import { CommonService } from '../../services/common.service';
import { LayoutService } from '../../services/layout.service';
/*MODELS*/
import { Transaction } from '../../models/transaction.model';
@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.scss']
})
export class TransactionComponent implements OnInit {

  public toggle: boolean;
  private _txHash: string;
  transaction: Observable<Transaction>;

  constructor(private _commonService: CommonService, private _route: ActivatedRoute, private _layoutService: LayoutService) {
  }
  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
    this.transaction = this._route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this._txHash = params.get('id');
        return this._commonService.getTransaction(this._txHash).pipe(
          tap(() => {
            this._layoutService.isPageLoading.next(false);
          })
        );
      })
    );
  }


}

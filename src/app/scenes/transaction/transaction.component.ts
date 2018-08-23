/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {Transaction} from '../../models/transaction.model';
import {CommonService} from '../../services/common.service';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {

  private _txHash: string;
  transaction: Observable<Transaction>;

  constructor(private _commonService: CommonService, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.transaction = this.route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this._txHash = params.get('id');
        return this._commonService.getTransaction(this._txHash);
      })
    );
  }

}

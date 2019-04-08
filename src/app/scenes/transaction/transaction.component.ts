/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {interval, Observable} from 'rxjs';
import {mergeMap, startWith, switchMap, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {WalletService} from '../../modules/wallet/wallet.service';
/*MODELS*/
import {Transaction} from '../../models/transaction.model';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.scss']
})
export class TransactionComponent implements OnInit, OnDestroy {

  showUtf8 = false;
  private _txHash: string;
  transaction$: Observable<Transaction> = this._route.paramMap.pipe(
    switchMap((params: ParamMap) => {
      this._txHash = params.get('id');
      return this._commonService.getTransaction(this._txHash).pipe(
        tap(() => {
          this._layoutService.offLoading();
        })
      );
    })
  );
  recentBlockNumber$: Observable<number> = interval(5000).pipe(
    startWith(0),
    mergeMap(() => fromPromise(this._walletService.w3.eth.getBlockNumber())),
    tap(res => console.log(res)),
  );

  constructor(private _commonService: CommonService,
              private _route: ActivatedRoute,
              private _layoutService: LayoutService,
              private _walletService: WalletService,
  ) {
  }

  async ngOnInit() {
    this._layoutService.onLoading();
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }
}

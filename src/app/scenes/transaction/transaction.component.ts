/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {interval, Observable, of, Subscription} from 'rxjs';
import {map, mergeMap, startWith, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {WalletService} from '../../modules/wallet/wallet.service';
/*MODELS*/
import {Transaction} from '../../models/transaction.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class TransactionComponent implements OnInit, OnDestroy {

  showUtf8 = false;
  tx: Transaction;

  recentBlockNumber$: Observable<number> = interval(5000).pipe(
    startWith(0),
    mergeMap(() => fromPromise(this._walletService.w3.eth.getBlockNumber())),
  );

  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService,
              private _route: ActivatedRoute,
              private _layoutService: LayoutService,
              private _walletService: WalletService,
  ) {
  }

  async ngOnInit() {
    this._layoutService.onLoading();
    this._subsArr$.push(
      this._route.paramMap.pipe(
        tap(() => {
          this._layoutService.onLoading();
        }),
        map((params: ParamMap) => params.get('id')),
        mergeMap((txHash: string) => this.getTx(txHash)),
      ).subscribe((tx: (Transaction | null)) => {
        this.tx = tx;
        this._layoutService.offLoading();
      })
    );
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  /**
   * getting tx from server
   * @param txHash
   */
  private getTx(txHash: string): Observable<Transaction | null> {
    return this._commonService.getTransaction(txHash).pipe(
      mergeMap((tx: Transaction | null) => {
        if (!tx) {
          return this._walletService.getTxData(txHash);
        }
        return of(tx);
      }),
    );
  }
}

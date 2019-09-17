/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {interval, Observable, of, Subscription} from 'rxjs';
import {map, mergeMap, startWith, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {WalletService} from '../../services/wallet.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {Transaction} from '../../models/transaction.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES} from '../../utils/constants';

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
    mergeMap(() => fromPromise<number>(this._walletService.w3.eth.getBlockNumber())),
  );

  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService,
              private _route: ActivatedRoute,
              private _layoutService: LayoutService,
              private _walletService: WalletService,
              private _metaService: MetaService,
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
        tx.logs = JSON.stringify(JSON.parse(tx.logs), null, '\t');
        this.tx = tx;
        this._layoutService.offLoading();
      })
    );
    this._metaService.setTitle(META_TITLES.TRANSACTION.title);
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

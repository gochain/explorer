/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {forkJoin, interval, Observable, of, Subscription} from 'rxjs';
import {filter, map, mergeMap, startWith, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {WalletService} from '../../modules/wallet/wallet.service';
/*MODELS*/
import {Transaction} from '../../models/transaction.model';
import {Transaction as Web3Tx} from 'web3/eth/types';
import {TransactionReceipt} from 'web3/types';
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
          return this.getPendingTx(txHash);
        }
        return of(tx);
      }),
    );
  }

  /**
   * getting tx from node in case of server haven't processed yet
   * @param txHash
   */
  private getPendingTx(txHash: string): Observable<Transaction | null> {
    return forkJoin(
      fromPromise(this._walletService.w3.eth.getTransaction(txHash)),
      fromPromise(this._walletService.w3.eth.getTransactionReceipt(txHash)),
    ).pipe(
      filter((res: [Web3Tx, TransactionReceipt]) => !!res[0]),
      map((res: [Web3Tx, TransactionReceipt]) => {
        const tx: Web3Tx = res[0];
        const txReceipt = res[1];
        const finalTx: Transaction = new Transaction();
        finalTx.tx_hash = tx.hash;
        finalTx.value = tx.value;
        finalTx.gas_price = tx.gasPrice;
        finalTx.gas_limit = '' + tx.gas;
        finalTx.nonce = '' + tx.nonce;
        finalTx.input_data = tx.input.replace(/^0x/, '');
        finalTx.from = tx.from;
        finalTx.to = tx.to;
        if (txReceipt) {
          finalTx.block_number = tx.blockNumber;
          finalTx.gas_fee = '' + (+tx.gasPrice * txReceipt.gasUsed);
          finalTx.contract_address = txReceipt.contractAddress;
          finalTx.status = txReceipt.status;
          finalTx.created_at = new Date();
        }
        return finalTx;
      }),
    );
  }
}

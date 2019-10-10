/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {forkJoin, interval, Observable, of, Subscription} from 'rxjs';
import {map, mergeMap, startWith, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {WalletService} from '../../services/wallet.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {ProcessedLog, Transaction, TxLog} from '../../models/transaction.model';
import {ContractEventsAbi} from '../../utils/types';
import {AbiItem} from 'web3-utils';
import {Address} from '../../models/address.model';
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
  showLogsRaw = false;
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

  ngOnInit() {
    this._layoutService.onLoading();
    this._subsArr$.push(
      this._route.paramMap.pipe(
        tap(() => {
          this._layoutService.onLoading();
        }),
        map((params: ParamMap) => params.get('id')),
        mergeMap((txHash: string) => this.getTx(txHash)),
      ).subscribe((tx: (Transaction | null)) => {
        tx.input_data = '0x' + tx.input_data;
        tx.parsedLogs = JSON.parse(tx.logs);
        tx.prettifiedLogs = JSON.stringify(tx.parsedLogs, null, '\t');
        this.tx = tx;
        this.processTransaction();
        this._layoutService.offLoading();
      })
    );
    this._metaService.setTitle(META_TITLES.TRANSACTION.title);
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  processTransaction() {
    if (!this.tx.parsedLogs.length) {
      return;
    }
    forkJoin([
      this._commonService.eventsAbi$,
      this._commonService.getAddress(this.tx.to),
    ]).subscribe(([events, address]: [ContractEventsAbi, Address]) => {
      this.tx.addressData = address;
      this.tx.processedLogs = this.tx.parsedLogs.map((log: TxLog) => {
        const prettyLog: ProcessedLog = new ProcessedLog();
        prettyLog.index = +log.logIndex;
        prettyLog.contract_address = log.address;
        prettyLog.removed = log.removed;
        let abiItem: AbiItem;
        const eventSignature = <string>log.topics[0];
        const knownEvent = events[eventSignature];
        if (knownEvent) {
          const ercType = address.erc_types.find((item: string) => !!knownEvent[item]);
          if (ercType) {
            abiItem = knownEvent[ercType];
          }
        }
        if (abiItem) {
          log.topics.shift();
          const decoded = this._walletService.w3.eth.abi.decodeLog(
            abiItem.inputs,
            log.data.replace('0x', ''),
            <string[]>log.topics
          );
          const res: string[] = abiItem.inputs.map(input => `${input.name}: ${decoded[input.name]}`);
          prettyLog.data = `${abiItem.name}(${res.join(', ')})`;
        } else {
          prettyLog.data = `topics: ${log.topics.join(',')} data: ${log.data}`;
        }
        return prettyLog;
      });
    });
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

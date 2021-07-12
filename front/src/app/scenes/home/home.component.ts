/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {interval, Observable} from 'rxjs';
import {mergeMap, startWith, tap} from 'rxjs/operators';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
import {CommonService} from '../../services/common.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {BlockList} from '../../models/block_list.model';
import {Stats} from '../../models/stats.model';
import {SupplyStats} from '../../models/supply.model';
import {META_TITLES} from '../../utils/constants';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})

export class HomeComponent implements OnInit, OnDestroy {
  stats$: Observable<Stats> = interval(300000).pipe(
    startWith(0),
    mergeMap(() => this._commonService.getStats())
  );

  private _supplyStatsInProgress = false;

  supplyStats$: Observable<SupplyStats> = interval(5000).pipe(
      startWith(0),
      mergeMap(async () => {
        if (this._supplyStatsInProgress) return; // don't stack
        try {
          this._supplyStatsInProgress = true;
          return await this._commonService.getSupplyStats().toPromise();
        } finally {
          this._supplyStatsInProgress = false;
        }
      })
  );

  private _recentBlocksInProgress = false;

  recentBlocks$: Observable<BlockList> = interval(5000).pipe(
    startWith(0),
    tap(() => {
      this._layoutService.offLoading();
    }),
    mergeMap(async () => {
      if (this._recentBlocksInProgress) return; // don't stack
      try {
        this._recentBlocksInProgress = true;
        return await this._commonService.getRecentBlocks().toPromise();
      } finally {
        this._recentBlocksInProgress = false;
      }
    }),
  );

  constructor(
    private _commonService: CommonService,
    private _layoutService: LayoutService,
    private _metaService: MetaService,
  ) {
  }

  ngOnInit() {
    this._layoutService.onLoading();
    this._metaService.setTitle(META_TITLES.HOME.title);
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }
}

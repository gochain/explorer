/*CORE*/
import {Component, OnInit} from '@angular/core';
import {interval, Observable} from 'rxjs';
import {mergeMap, startWith, tap} from 'rxjs/operators';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {BlockList} from '../../models/block_list.model';
import {Stats} from '../../models/stats.model';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  stats$: Observable<Stats> = interval(120000).pipe(
    startWith(0),
    mergeMap(() => this._commonService.getStats())
  );

  recentBlocks$: Observable<BlockList> = interval(5000).pipe(
    startWith(0),
    tap(() => {
      this._layoutService.isPageLoading.next(false);
    }),
    mergeMap(() => this._commonService.getRecentBlocks()),
  );

  constructor(private _commonService: CommonService, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
  }
}

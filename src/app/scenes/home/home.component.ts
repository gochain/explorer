/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {interval, Observable, Subscription} from 'rxjs';
import {startWith} from 'rxjs/operators';
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
export class HomeComponent implements OnInit, OnDestroy {

  recentBlocks: BlockList;
  stats$: Observable<Stats>;
  private _sub: Subscription;

  constructor(private _commonService: CommonService, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
    // to-do: replace to ws
    this._sub = interval(5000).pipe(
      startWith(0)
    ).subscribe(() => {
      this._commonService.getRecentBlocks().subscribe((data: BlockList) => {
          this.recentBlocks = data;
          this._layoutService.isPageLoading.next(false);
        }
      );
    });

    this.stats$ = this._commonService.getStats();
  }

  ngOnDestroy() {
    this._sub.unsubscribe();
  }
}

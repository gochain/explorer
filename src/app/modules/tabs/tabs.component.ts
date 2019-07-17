/*CORE*/
import {AfterContentInit, Component, ContentChildren, EventEmitter, Input, OnInit, Output, QueryList} from '@angular/core';
import {ActivatedRoute, Params, Router} from '@angular/router';
import {Subscription} from 'rxjs';
/*COMPONENTS*/
import {TabComponent} from './components/tab/tab.component';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-tabs',
  templateUrl: 'tabs.component.html',
  styleUrls: ['./tabs.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class TabsComponent implements OnInit, AfterContentInit {
  @Input() name: string;
  @Input() disabled = false;
  @Output() changeEmitter = new EventEmitter<any>();
  @ContentChildren(TabComponent) tabs: QueryList<TabComponent>;
  activeTab: TabComponent;

  private _tabName: string;
  private _subsArr$: Subscription[] = [];

  constructor(private _activatedRoute: ActivatedRoute, private _router: Router) {
  }

  ngOnInit() {
    // setting initial tab if query param provided
    this._subsArr$.push(this._activatedRoute.queryParams.subscribe((params: Params) => {
      this._tabName = params[this.name] || null;
    }));
  }

  ngAfterContentInit(): void {
    // subscribing for tabs change if tabs will be added or deleted
    this._subsArr$.push(this.tabs.changes.subscribe(() => this.onTabsChange()));
    // asynchronous update preventing ExpressionChangedAfterItHasBeenCheckedError
    setTimeout(() => {
      this.findTab();
    });
  }

  findTab(): void {
    let activeTab: TabComponent;
    if (this._tabName) {
      activeTab = this.tabs.find((tab: TabComponent) => tab.name === this._tabName) || this.tabs.first;
    } else {
      activeTab = this.tabs.first;
    }
    this.onTabSelect(activeTab, false);
  }

  onTabsChange(): void {
    if (this.tabs.length) {
      this.findTab();
    } else {
      this.activeTab = null;
    }
  }

  /**
   * replacing url so query params won't affect url history
   * @param tab
   * @param emit
   * @param replaceUrl
   */
  onTabSelect(tab: TabComponent, emit = true, replaceUrl = true) {
    if (this.disabled) {
      return;
    }
    if (emit && this.changeEmitter) {
      this.changeEmitter.emit(this.activeTab.name);
    }
    if (this.activeTab) {
      this.activeTab.content.active = false;
    }
    this.activeTab = tab;
    this.activeTab.content.active = true;
    this._router.navigate([], {
      relativeTo: this._activatedRoute,
      queryParams: {
        ...this._activatedRoute.snapshot.queryParams,
        [this.name]: tab.name,
      },
      replaceUrl,
    });
  }
}

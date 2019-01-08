/*CORE*/
import {AfterContentInit, Component, ContentChildren, EventEmitter, Input, OnInit, Output, QueryList} from '@angular/core';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
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
  @Output() onChangeEmitter = new EventEmitter<void>();
  @ContentChildren(TabComponent) tabs: QueryList<TabComponent>;
  activeTab: TabComponent;

  private _initialTabName: string;
  private _subsArr$: Subscription[] = [];

  constructor(private _activatedRoute: ActivatedRoute, private _router: Router) {
  }

  ngOnInit() {
    this._subsArr$.push(this._activatedRoute.queryParamMap.subscribe((params: ParamMap) => {
      if (params.has(this.name)) {
        this._initialTabName = params.get(this.name);
      }
    }));
  }

  ngAfterContentInit() {
    this._subsArr$.push(this.tabs.changes.subscribe(this.onTabsChange));
    // Asynchronous update preventing ExpressionChangedAfterItHasBeenCheckedError
    setTimeout(() => {
      if (this._initialTabName) {
        const activeTab = this.tabs.find((tab: TabComponent) => tab.name === this._initialTabName) || this.tabs.first;
        this.onTabSelect(activeTab);
      } else {
        this.onTabSelect(this.tabs.first);
      }
    });
  }

  onTabsChange = () => {
    if (this.tabs.length) {
      const exist = this.tabs.some(tab => tab === this.activeTab);
      if (!exist) {
        this.activeTab = this.tabs.first;
      }
    } else {
      this.activeTab = null;
    }
  }

  onTabSelect(tab: TabComponent) {
    if (this.onChangeEmitter) {
      this.onChangeEmitter.emit();
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
      }
    });
  }
}

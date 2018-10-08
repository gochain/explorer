import {AfterContentInit, Component, ContentChildren, QueryList} from '@angular/core';
import {TabComponent} from './components/tab/tab.component';
import {Subscription} from 'rxjs';
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-tabs',
  templateUrl: 'tabs.component.html',
  styleUrls: ['./tabs.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class TabsComponent implements AfterContentInit {
  @ContentChildren(TabComponent) tabs: QueryList<TabComponent>;
  activeTab: TabComponent;
  private _subsArr$: Subscription[] = [];

  ngAfterContentInit() {
    this._subsArr$.push(this.tabs.changes.subscribe(this.onTabsChange));
    this.activeTab = this.tabs.first;
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
  };

  onTabSelect(tab: TabComponent) {
    this.activeTab = tab;
  }
}

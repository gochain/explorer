import { AfterContentInit, Component, ContentChildren, QueryList } from '@angular/core';
import { TabComponent } from './components/tab/tab.component';
@Component({
  selector: 'app-tabs',
  templateUrl: 'tabs.component.html',
  styleUrls: ['./tabs.component.scss']
})
export class TabsComponent implements AfterContentInit {
  @ContentChildren(TabComponent) tabs: QueryList<TabComponent>;
  activeTab: TabComponent;

  ngAfterContentInit() {
    this.tabs.changes.subscribe(this.onTabsChange);
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
  }

  onTabSelect(tab: TabComponent) {
    this.activeTab = tab;
  }
}

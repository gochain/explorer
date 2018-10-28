import {Component, ContentChild, Input} from '@angular/core';
import { TabTitleComponent } from '../tab-title/tab-title.component';
import { TabContentComponent } from '../tab-content/tab-content.component';

@Component({
  selector: 'tab',
  templateUrl: 'tab.component.html'
})
export class TabComponent {
  @Input() name: string;
  @ContentChild(TabTitleComponent) title: TabTitleComponent;
  @ContentChild(TabContentComponent) content: TabContentComponent;
}
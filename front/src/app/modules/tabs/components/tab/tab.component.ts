import {Component, ContentChild, Input} from '@angular/core';
import {TabContentComponent} from '../tab-content/tab-content.component';

@Component({
  selector: 'tab',
  templateUrl: 'tab.component.html'
})
export class TabComponent {
  @Input() name: string;
  @Input() title: string;
  @Input() description: string;
  /*@ContentChild(TabTitleComponent) title: TabTitleComponent;*/
  @ContentChild(TabContentComponent, {static: true}) content: TabContentComponent;
}

import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TabsComponent } from './tabs.component';
import { TabTitleComponent } from './components/tab-title/tab-title.component';
import { TabContentComponent } from './components/tab-content/tab-content.component';
import { TabComponent } from './components/tab/tab.component';
import {TooltipModule} from 'ng2-tooltip-directive';

@NgModule({
  declarations: [TabsComponent, TabTitleComponent, TabContentComponent, TabComponent],
  imports: [CommonModule, TooltipModule],
  exports: [TabsComponent, TabTitleComponent, TabContentComponent, TabComponent],
  providers: []
})
export class TabsModule {}

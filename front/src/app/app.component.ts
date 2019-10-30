import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {LayoutService} from './services/layout.service';
import { AutoUnsubscribe } from './decorators/auto-unsubscribe';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class AppComponent {
  constructor(private _layoutService: LayoutService) {
  }
}

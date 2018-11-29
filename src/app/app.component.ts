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
export class AppComponent implements OnInit {
  isPageLoading = false;

  private _subsArr$: Subscription[] = [];

  constructor(private _layoutService: LayoutService, private _cdr: ChangeDetectorRef) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._layoutService.isPageLoading.subscribe((state: boolean) => {
        this.isPageLoading = state;
        this._cdr.detectChanges();
      })
    );
  }
}

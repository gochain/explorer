import {ChangeDetectorRef, Component, OnInit} from '@angular/core';
import {LayoutService} from './services/layout.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  isPageLoading = false;

  constructor(private _layoutService: LayoutService, private _cdr: ChangeDetectorRef) {
  }

  ngOnInit() {
    this._layoutService.isPageLoading.subscribe((state: boolean) => {
      this.isPageLoading = state;
      this._cdr.detectChanges();
    });
  }
}

import {Component, OnInit} from '@angular/core';
import {LayoutService} from './services/template.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  isSidenavOpen = true;
  isPageLoading = false;

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isSidenavOpen.subscribe((state: boolean) => {
      this.isSidenavOpen = state;
    });
    this._layoutService.isPageLoading.subscribe((state: boolean) => {
      this.isPageLoading = state;
    });
  }
}

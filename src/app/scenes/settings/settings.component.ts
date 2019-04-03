/*
import {Component, OnInit} from '@angular/core';
import {LayoutService} from '../../services/layout.service';
import {filter} from 'rxjs/operators';

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss']
})
export class SettingsComponent implements OnInit {
  themeColor$: string;

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.themeColor$.pipe(
      filter(value => !!value)
    ).subscribe(value => {
      this.themeColor$ = value;
    });
  }

  onThemeColorChange() {
    this._layoutService.themeColor$.next(this.themeColor$);
  }

  onSubmit() {

  }
}
*/

/*
import {Component, OnInit} from '@angular/core';
import {filter} from 'rxjs/operators';
import {LayoutService} from '../../services/layout.service';
import {ThemeColor} from '../../utils/enums';

@Component({
  selector: 'app-toggle-switch',
  templateUrl: './toggle-switch.component.html',
  styleUrls: ['./toggle-switch.component.scss']
})
export class ToggleSwitchComponent implements OnInit {
  themeColor: ThemeColor;

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.themeColor$.pipe(
      filter(value => !!value)
    ).subscribe(value => {
      this.themeColor = value;
    });
  }

  onChange() {
    const color = this.themeColor === ThemeColor.DARK ? ThemeColor.LIGHT : ThemeColor.DARK;
    this._layoutService.themeColor$.next(color);
  }
}
*/

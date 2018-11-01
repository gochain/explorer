import {Component, OnInit} from '@angular/core';
import {filter} from 'rxjs/operators';
import {LayoutService} from '../../services/layout.service';

@Component({
  selector: 'app-toggle-switch',
  templateUrl: './toggle-switch.component.html',
  styleUrls: ['./toggle-switch.component.scss']
})
export class ToggleSwitchComponent implements OnInit {
  themeColor: string;

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.themeColor.pipe(
      filter(value => !!value)
    ).subscribe(value => {
      this.themeColor = value;
    });
  }

  onChange() {
    const color = this.themeColor === 'dark' ? 'light' : 'dark';
    this._layoutService.themeColor.next(color);
  }
}

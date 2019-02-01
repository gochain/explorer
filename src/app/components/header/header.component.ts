/*CORE*/
import {Component, OnInit} from '@angular/core';
/*MODELS*/
import {MenuItem} from '../../models/menu_item.model';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
/*UTILS*/
import {MENU_ITEMS} from '../../utils/constants';
import {ThemeColor} from '../../utils/enums';

const LOGO_NAME = {
  [ThemeColor.LIGHT]: 'logo_fullcolor.png',
  [ThemeColor.DARK]: 'logo_allwhite.png',
}

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  themeColor: ThemeColor;
  logoSrc = LOGO_NAME[ThemeColor.DARK];
  navItems: MenuItem[] = MENU_ITEMS;

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.themeColor.subscribe((value: ThemeColor) => {
      this.themeColor = value;
      this.logoSrc = LOGO_NAME[value];
    });
  }
}

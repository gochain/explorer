/*CORE*/
import {Component, OnInit} from '@angular/core';
/*MODELS*/
import {MenuItem} from '../../models/menu_item.model';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
/*UTILS*/
import {MENU_ITEMS} from '../../utils/constants';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  themeColor: string;

  navItems: MenuItem[] = MENU_ITEMS;

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.themeColor.subscribe(value => {
      this.themeColor = value;
    });
  }
}

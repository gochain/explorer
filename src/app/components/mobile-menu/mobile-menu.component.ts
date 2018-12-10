/*CORE*/
import {Component, OnInit} from '@angular/core';
/*MODELS*/
import {MenuItem} from '../../models/menu_item.model';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
/*UTILS*/
import {MENU_ITEMS} from '../../utils/constants';

@Component({
  selector: 'app-mobile-menu',
  templateUrl: './mobile-menu.component.html',
  styleUrls: ['./mobile-menu.component.scss']
})
export class MobileMenuComponent implements OnInit {

  navItems: MenuItem[] = MENU_ITEMS;

  constructor(public layoutService: LayoutService) {
  }

  ngOnInit() {
  }

  hideMenu() {
    this.layoutService.mobileMenuState.next(false);
  }
}

/*CORE*/
import {Component} from '@angular/core';
/*MODELS*/
import {MenuItem} from '../../models/menu_item.model';
/*SERVICES*/
/*UTILS*/
import {MENU_ITEMS} from '../../utils/constants';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
  navItems: MenuItem[] = MENU_ITEMS;

  constructor() {
  }
}

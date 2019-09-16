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
  showMenu = false;
  showSearch = false;
  navItems: MenuItem[] = MENU_ITEMS;

  constructor() {
  }

  toggleMenu(): void {
    this.showMenu = !this.showMenu;
    this.showSearch = false;
  }

  toggleSearch(): void {
    this.showSearch = !this.showSearch;
    this.showMenu = false;
  }
}

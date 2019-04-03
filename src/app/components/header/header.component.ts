/*CORE*/
import {Component} from '@angular/core';
import {map} from 'rxjs/operators';
/*MODELS*/
import {MenuItem} from '../../models/menu_item.model';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
/*UTILS*/
import {LOGO_NAMES, MENU_ITEMS} from '../../utils/constants';
import {ThemeColor} from '../../utils/enums';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
  logoSrc$ = this._layoutService.themeColor$.pipe(
    map((value: ThemeColor) => LOGO_NAMES[value]),
  );
  navItems: MenuItem[] = MENU_ITEMS;

  constructor(private _layoutService: LayoutService) {
  }
}

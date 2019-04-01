/*CORE*/
import {Component} from '@angular/core';
import {map} from 'rxjs/operators';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
/*UTILS*/
import {ThemeColor} from '../../utils/enums';
import {LOGO_NAMES} from '../../utils/constants';

@Component({
  selector: 'app-mobile-header',
  templateUrl: './mobile-header.component.html',
  styleUrls: ['./mobile-header.component.scss']
})
export class MobileHeaderComponent {

  themeColor: string;
  logoSrc$ = this.layoutService.themeColor$.pipe(
    map((value: ThemeColor) => LOGO_NAMES[value]),
  );

  constructor(public layoutService: LayoutService) {
  }

  toggleMenu() {
    this.layoutService.mobileMenuState.next(!this.layoutService.mobileMenuState.value);
    this.layoutService.mobileSearchState.next(false);
  }

  toggleSearch() {
    this.layoutService.mobileSearchState.next(!this.layoutService.mobileSearchState.value);
    this.layoutService.mobileMenuState.next(false);
  }
}

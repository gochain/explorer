import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {ROUTES} from '../../utils/routes';
import {LayoutService} from '../../services/layout.service';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
})
export class SearchComponent {
  value = '';

  constructor(private router: Router, public layoutService: LayoutService) {
  }

  async search() {
    if (this.value.length === 42) {
      this.layoutService.mobileSearchState.next(false);
      await this.router.navigate([`/${ROUTES.ADDRESS}/`, this.value]);
    } else if (this.value.length === 66) {
      this.layoutService.mobileSearchState.next(false);
      await this.router.navigate([`/${ROUTES.TRANSACTION}/`, this.value]);
    }
  }
}

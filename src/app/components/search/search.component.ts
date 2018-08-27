import {Component} from '@angular/core';
import {Router} from '@angular/router';
import {ROUTES} from '../../utils/routes';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
})
export class SearchComponent {
  value = '';

  constructor(private router: Router) {
  }

  async search() {
    if (this.value.length === 42) {
      await this.router.navigate([`/${ROUTES.ADDRESS}/`, this.value]);
    } else if (this.value.length === 66) {
      await this.router.navigate([`/${ROUTES.TRANSACTION}/`, this.value]);
    }
  }
}

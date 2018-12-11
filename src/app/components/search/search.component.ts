/*CORE*/
import {Component} from '@angular/core';
import {Router} from '@angular/router';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
/*UTILS*/
import {ROUTES} from '../../utils/constants';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
})
export class SearchComponent {
  value = '';

  constructor(private router: Router, public layoutService: LayoutService, private toastrService: ToastrService) {
  }

  async search() {
    if (this.value.length === 42) {
      this.layoutService.mobileSearchState.next(false);
      await this.router.navigate([`/${ROUTES.ADDRESS}/`, this.value]);
    } else if (this.value.length === 66) {
      this.layoutService.mobileSearchState.next(false);
      await this.router.navigate([`/${ROUTES.TRANSACTION}/`, this.value]);
    } else {
      this.toastrService.warning('the data you entered is not valid');
    }
  }
}

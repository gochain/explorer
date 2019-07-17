/*CORE*/
import {Component} from '@angular/core';
import {Router} from '@angular/router';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
import {CommonService} from '../../services/common.service';
/*UTILS*/
import {ROUTES} from '../../utils/constants';

@Component({
  selector: 'app-search',
  templateUrl: './search.component.html',
  styleUrls: ['./search.component.scss'],
})
export class SearchComponent {
  value = '';

  constructor(
    private router: Router,
    public layoutService: LayoutService,
    private toastrService: ToastrService,
    private commonService: CommonService,
  ) {
  }

  async search() {
    const value = this.value.trim();
    if (value.length === 42) {
      await this.router.navigate([`/${ROUTES.ADDRESS}/`, value]);
    } else if (value.length === 66) {
      this.commonService.checkBlockExist(value).subscribe(resp =>  {
        this.router.navigate([`/${ROUTES.BLOCK}/`, value]);
      }, error => {
        this.commonService.checkTransactionExist(value).subscribe(resp =>  {
          this.router.navigate([`/${ROUTES.TRANSACTION}/`, value]);
          }, error => {
            this.toastrService.warning('the data you entered is not valid');
        });
      });
    } else if (value.length < 8 && /^\d+$/.test(value)) {
      await this.router.navigate([`/${ROUTES.BLOCK}/`, value]);
    } else {
      this.toastrService.warning('the data you entered is not valid');
      return;
    }

    this.layoutService.mobileSearchState.next(false);
  }
}

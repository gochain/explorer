/*CORE*/
import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable()
export class LayoutService {
  isSidenavOpen: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);
  isPageLoading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  toggleSidenav() {
    this.isSidenavOpen.next(!this.isSidenavOpen.value);
  }

  toggleLoading() {
    this.isPageLoading.next(!this.isPageLoading.value);
  }
}

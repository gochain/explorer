/*CORE*/
import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable()
export class LayoutService {
  isSidenavOpen: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(true);

  toggleSidenav() {
    this.isSidenavOpen.next(!this.isSidenavOpen.value);
  }
}

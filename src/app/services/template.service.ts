/*CORE*/
import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';

@Injectable()
export class LayoutService {
  isPageLoading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  toggleLoading() {
    this.isPageLoading.next(!this.isPageLoading.value);
  }
}

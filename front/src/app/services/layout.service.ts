/*CORE*/
import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
/*MODELS*/
import {ThemeSettings} from '../models/theme_settings.model';
/*UTILS*/
import {THEME_SETTINGS} from '../utils/constants';
import {ThemeColor} from '../utils/enums';

@Injectable({
  providedIn: 'root',
})
export class LayoutService {
  isPageLoading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  themeColor$: BehaviorSubject<ThemeColor> = new BehaviorSubject<ThemeColor>(ThemeColor.LIGHT);
  themeSettings: ThemeSettings;
  mobileMenuState: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  mobileSearchState: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  constructor() {
    const themeSettings = localStorage.getItem('THEME_SETTINGS');
    if (themeSettings) {
      this.themeSettings = JSON.parse(themeSettings);
    } else {
      localStorage.setItem('THEME_SETTINGS', JSON.stringify(THEME_SETTINGS));
      this.themeSettings = THEME_SETTINGS;
    }

    this.themeColor$.next(this.themeSettings.color);

    this.themeColor$.subscribe((value: ThemeColor) => {
      // document.body.classList.remove(ThemeColor.DARK, ThemeColor.LIGHT);
      // document.body.classList.add(value);
      // this.themeSettings.color = value;
      // localStorage.setItem('THEME_SETTINGS', JSON.stringify(this.themeSettings));
    });
  }

  toggleLoading() {
    this.isPageLoading.next(!this.isPageLoading.value);
  }

  onLoading() {
    this.isPageLoading.next(true);
  }

  offLoading() {
    this.isPageLoading.next(false);
  }
}

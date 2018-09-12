/*CORE*/
import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
/*MODELS*/
import {ThemeSettings} from '../models/theme_settings.model';
/*UTILS*/
import {THEME_SETTINGS} from '../utils/constants';

@Injectable()
export class LayoutService {
  isPageLoading: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);
  themeColor: BehaviorSubject<string>;
  themeSettings: ThemeSettings;

  constructor() {
    const themeSettings = localStorage.getItem('THEME_SETTINGS');
    if (themeSettings) {
      this.themeSettings = JSON.parse(themeSettings);
    } else {
      localStorage.setItem('THEME_SETTINGS', JSON.stringify(THEME_SETTINGS));
      this.themeSettings = THEME_SETTINGS;
    }
    this.themeColor = new BehaviorSubject<string>(this.themeSettings.color);

    this.themeColor.subscribe(value => {
      document.body.classList.remove('dark', 'light');
      document.body.classList.add(value);
      this.themeSettings.color = value;
      localStorage.setItem('THEME_SETTINGS', JSON.stringify(this.themeSettings));
    });
  }

  toggleLoading() {
    this.isPageLoading.next(!this.isPageLoading.value);
  }
}
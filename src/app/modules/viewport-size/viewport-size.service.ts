import {Inject, Injectable, NgZone} from '@angular/core';
import {ViewportSizeEnum} from './viewport-size.enum';
import {IConfig} from './config.interface';
import {BehaviorSubject, fromEvent} from 'rxjs';
import {debounceTime, distinctUntilChanged} from 'rxjs/operators';

@Injectable()
export class ViewportSizeService {
  size$: BehaviorSubject<ViewportSizeEnum> = new BehaviorSubject<ViewportSizeEnum>(null);

  constructor(@Inject('config') private _config: IConfig,
              private _zone: NgZone
  ) {
    this.onWindowResize(window.innerWidth);
    this.initTracker();
  }

  initTracker() {
    this._zone.runOutsideAngular(() => {
      fromEvent(window, 'resize').pipe(
        debounceTime(500),
        distinctUntilChanged()
      ).subscribe((e: any) => {
        this._zone.run(() => {
          this.onWindowResize(e.target.innerWidth);
        });
      });
    });
  }

  onWindowResize = (windowWidth: number) => {
    let currentSize: ViewportSizeEnum;

    if (windowWidth >= this._config.large) {
      currentSize = ViewportSizeEnum.large;
    } else if (windowWidth >= this._config.medium) {
      currentSize = ViewportSizeEnum.medium;
    } else {
      currentSize = ViewportSizeEnum.small;
    }
    this.size$.next(currentSize);
  }
}

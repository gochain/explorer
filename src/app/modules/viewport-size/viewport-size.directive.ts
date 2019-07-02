/*CORE*/
import {Directive, Input, OnDestroy, OnInit, TemplateRef, ViewContainerRef} from '@angular/core';
import {Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
/*SERVICES*/
import {ViewportSizeService} from './viewport-size.service';
/*UTILS*/
import {ViewportSizeEnum} from './viewport-size.enum';


@Directive({selector: '[appViewportSize]'})
export class ViewportSizeDirective implements OnInit, OnDestroy {
  private _visibleSize: ViewportSizeEnum[];
  private _embedded = false;
  private _sub: Subscription;

  constructor(private _viewportSizeService: ViewportSizeService,
              private _templateRef: TemplateRef<any>,
              private _viewContainer: ViewContainerRef,
  ) {
  }

  @Input() set appViewportSize(sizes: ViewportSizeEnum[]) {
    this._visibleSize = sizes;
  }

  ngOnInit() {
    this._sub = this._viewportSizeService.size$
      .pipe(
        filter(currentSize => currentSize !== null)
      )
      .subscribe((currentSize: ViewportSizeEnum) => {
        this.onResize(currentSize);
      });
  }

  ngOnDestroy() {
    this._sub.unsubscribe();
  }

  onResize(currentSize: ViewportSizeEnum) {
    if (this._visibleSize.includes(currentSize)) {
      if (!this._embedded) {
        this._embedded = true;
        this._viewContainer.createEmbeddedView(this._templateRef);
      }
    } else {
      this._embedded = false;
      this._viewContainer.clear();
    }
  }
}

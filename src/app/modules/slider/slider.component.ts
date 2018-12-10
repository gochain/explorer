/*CORE*/
import { Component, OnInit, Input } from '@angular/core';

export interface ISliderOptions {
    start: number;
    step: number;
    sensitivity: number;
}

@Component({
  selector: 'app-slider',
  templateUrl: './slider.component.html',
  styleUrls: ['./slider.component.scss']
})

export class SliderComponent implements OnInit {

  @Input() stats: any;
  @Input() titles: string[];
  @Input() options: ISliderOptions;
  private _initialPoint: any;
  private _finalPoint: any;
  private _touchOffsetX: number;
  private _slidePosition: number;

  constructor() { }

  ngOnInit(): void {
    this._slidePosition = this.options.start;
  }

  touchStart(event: TouchEvent): void {
    this._initialPoint = event.changedTouches[0];
    this._touchOffsetX = this._initialPoint.pageX - this._initialPoint.target.offsetLeft;
  }

  touchEnd(event: TouchEvent): void {
    this._finalPoint = event.changedTouches[0];
    if (this._finalPoint.pageX - this._initialPoint.pageX < -this.options.sensitivity) {
        if (-this.options.step !== this._slidePosition) {
            this._slidePosition = this._slidePosition - this.options.step;
        }
    } else if ((this._finalPoint.pageX - this._initialPoint.pageX) > this.options.sensitivity) {
        if (this.options.step !== this._slidePosition) {
            this._slidePosition = this._slidePosition + this.options.step;
        }
    }
  }
  onClickDot(index: number) {
    this._slidePosition = this.options.start - this.options.step * index;
  }
}

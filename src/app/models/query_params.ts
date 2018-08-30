import {BehaviorSubject, Observable, Subject} from 'rxjs';
import {HttpParams} from '@angular/common/http';

export class QueryParams {
  limit: number;
  skip: number;
  page: number;
  total: number;
  totalPage: number;
  currentTotal: number;
  state: Subject<number> = new Subject();
  private _state = 0;

  constructor(limit?: number) {
    this.limit = limit || 20;
    this.skip = 0;
    this.page = 1;
    this.currentTotal = this.limit;
  }

  next() {
    this.page++;
    this.skip += this.limit;
    this.currentTotal = this.page * this.limit;
    this.state.next(++this._state);
  }

  previous() {
    this.page--;
    this.skip -= this.limit;
    this.state.next(++this._state);
  }

  toPage(page: number) {
    this.page = page;
    this.skip = (this.page - 1) * this.limit;
    this.state.next(++this._state);
  }

  toStart() {
    this.page = 1;
    this.skip = 0;
    this.state.next(++this._state);
  }

  toEnd() {
    this.page = Math.ceil(this.total / this.limit);
    this.skip = (this.page - 1) * this.limit;
    this.state.next(++this._state);
  }

  get params() {
   return {limit: this.limit, skip: this.skip};
  }
}

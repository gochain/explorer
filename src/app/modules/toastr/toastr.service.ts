import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
import {Toastr} from './toastr.interface';

const CLOSE_TIME = 10000;

@Injectable({
  providedIn: 'root'
})
export class ToastrService {
  counter = 0;
  items$: BehaviorSubject<Toastr[]> = new BehaviorSubject<Toastr[]>([]);
  items: Toastr[] = [];

  constructor() {
  }

  success(msg: string) {
    this.add(msg, 'success');
  }

  warning(msg: string) {
    this.add(msg, 'warning');
  }

  danger(msg: string) {
    this.add(msg, 'danger');
  }

  info(msg: string) {
    this.add(msg, 'info');
  }

  add(msg: string, type: string) {
    const item: Toastr = {
      id: this.counter++,
      content: msg,
      type: type,
    };
    this.items = [...this.items, item];
    this.apply();
    setTimeout(() => this.delete(item.id), CLOSE_TIME);
  }

  delete(id: number) {
    this.items = this.items.filter((item: Toastr) => item.id !== id);
    this.apply();
  }

  apply() {
    this.items$.next(this.items);
  }
}

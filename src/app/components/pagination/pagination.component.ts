/*CORE*/
import {Component, EventEmitter, Input, Output} from '@angular/core';
import {QueryParams} from '../../models/query_params';

enum PaginationType {
  PAGE = 'page',
  DOTS = 'dot'
}

class PaginationItem {
  value: number | string;
  type: PaginationType;


  constructor(value: any, type: PaginationType) {
    this.value = value;
    this.type = type;
  }
}

@Component({
  selector: 'app-pagination',
  templateUrl: './pagination.component.html',
  styleUrls: ['./pagination.component.scss']
})
export class PaginationComponent {
  middle: number[] = [];
  current: number;

  _total: number;

  items: PaginationItem[] = [];

  perPage: number[] = [10, 25, 50, 100];

  _queryParam: QueryParams;

  @Input()
  set queryParam(queryParam: QueryParams) {
    this._queryParam = queryParam;
    this._queryParam.totalPage$.subscribe(value => {
      this._total = value;
      this.current = 1;
      this.calculateItems();
    });
    this._queryParam.page$.subscribe((page: number) => {
      this.current = page;
    });
  }

  @Output()
  onSelect = new EventEmitter<number>();

  onPageSelect(page: number) {
    if (page === this.current) {
      return;
    }
    this.current = page;
    this.calculateItems();
    this._queryParam.toPage(page);
    if (this.onSelect) {
      this.onSelect.emit(this.current);
    }
  }

  calculateItems() {
    const items = [];
    const dots: PaginationItem = new PaginationItem('...', PaginationType.DOTS);
    if (this._total < 12) {
      for (let i = 1; i <= this._total; i++) {
        const item: PaginationItem = new PaginationItem(i, PaginationType.PAGE);
        items.push(item);
      }
    } else if (this.current < 7) {
      for (let i = 1; i <= 10; i++) {
        const item: PaginationItem = new PaginationItem(i, PaginationType.PAGE);
        items.push(item);
      }
      const last: PaginationItem = new PaginationItem(this._total, PaginationType.PAGE);
      items.push(dots, last);
    } else if (this.current > (this._total - 6)) {
      const first: PaginationItem = new PaginationItem(1, PaginationType.PAGE);
      items.push(first, dots);

      for (let i = this._total - 10; i <= this._total; i++) {
        const item: PaginationItem = new PaginationItem(i, PaginationType.PAGE);
        items.push(item);
      }
    } else {
      const first: PaginationItem = new PaginationItem(1, PaginationType.PAGE);
      items.push(first, dots);

      for (let i = this.current - 4; i <= this.current + 4; i++) {
        const item: PaginationItem = new PaginationItem(i, PaginationType.PAGE);
        items.push(item);
      }

      const last: PaginationItem = new PaginationItem(this._total, PaginationType.PAGE);
      items.push(dots, last);
    }

    this.items = items;
  }

  next() {
    if (this.current !== this._total) {
      this.onPageSelect(this.current + 1);
    }
  }

  previous() {
    if (this.current !== 1) {
      this.onPageSelect(this.current - 1);
    }
  }
}

/*CORE*/
import {Injectable} from '@angular/core';
import {Title} from '@angular/platform-browser';

const DEFAULT_TITLE = 'GoChain Explorer';

@Injectable({
  providedIn: 'root'
})
export class MetaService {

  constructor(private title: Title) {
  }

  setTitle(value: string): void {
    this.title.setTitle(`${value} - ${DEFAULT_TITLE}`);
  }

  resetTitle(): void {
    this.title.setTitle(DEFAULT_TITLE);
  }
}

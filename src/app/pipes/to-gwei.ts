import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'toGwei'
})
export class ToGwei implements PipeTransform {

  transform(val: number): string {
    if (!val) {
      return '';
    }

    return Math.round(val / 1e7) / 1e2 + ' gwei';
  }
}

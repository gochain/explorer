import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'toGwei'
})
export class ToGwei implements PipeTransform {

  transform(val: number): string {
    if (!val) {
      return '';
    }

    return val / 1e9 + ' gwei';
  }
}

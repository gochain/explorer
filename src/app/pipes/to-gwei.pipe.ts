import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'toGwei'
})
export class ToGweiPipe implements PipeTransform {

  transform(val: number): string {
    if (!val) {
      return '';
    }

    return val / 1e9 + ' gwei';
  }
}

import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'bigNumber'
})

export class BigNumberPipe implements PipeTransform {

  transform(val: string): string {
    if (val == null) {
      return val;
    }
    const parts = val.toString().split('.');
    parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');
    return parts.join('.');
  }
}

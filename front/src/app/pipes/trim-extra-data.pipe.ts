import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'trimExtra'
})
export class TrimExtra implements PipeTransform {

  transform(val: string): string {
    if (!val) {
      return '';
    }
    return val.substring(0, val.indexOf('\u0000'));
  }
}

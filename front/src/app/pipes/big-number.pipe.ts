import {PipeTransform, Pipe} from '@angular/core';
import {numberWithCommas} from '../utils/functions';

@Pipe({
  name: 'bigNumber'
})

export class BigNumberPipe implements PipeTransform {

  transform(val: string): string {
    return numberWithCommas(val);
  }
}

import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'weiToGO'
})

export class WeiToGOPipe implements PipeTransform {

  transform(val: string, showUnit: boolean = true, fixedFraction: number = null): string {
    const moveTo = 18;
    const parts = val.toString().split('.');
    if (parts[0].length > moveTo) {
      parts[0] = parts[0].slice(0, parts[0].length - moveTo) + '.' + parts[0].slice(parts[0].length - moveTo, parts[0].length);
    } else {
      parts[0] = '0.' + '0'.repeat(moveTo - parts[0].length) + parts[0];
    }
    let value = parts.join('').toString();

    if (fixedFraction) {
      value = (+ value).toFixed(fixedFraction);
    }
    if (showUnit) {
      value += ' GO';
    }
    return value;
  }
}

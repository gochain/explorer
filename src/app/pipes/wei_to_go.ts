import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'weiToGO'
})

export class WeiToGOPipe implements PipeTransform {

  transform(val: string, showUnit: boolean = true, removeTrailingZeros: boolean = null): string {
    console.log(val);
    const moveTo = 18;
    const parts = val.toString().split('.');
    if (parts[0].length > moveTo) {
      parts[0] = parts[0].slice(0, parts[0].length - moveTo) + '.' + parts[0].slice(parts[0].length - moveTo, parts[0].length);
    } else {
      parts[0] = '0.' + '0'.repeat(moveTo - parts[0].length) + parts[0];
    }
    let value: string = parts.join('').toString();

    if (removeTrailingZeros) {
      // replace trailing zeros with exact amount of spaces
      value = value.replace(/0(?=(0+$|$))/g, ` `);
      value = value.replace(/\.(?=\s)/g, ` `);
    } else {
      // delete trailing zeros
      value = value.replace(/\.?0+$/, '');
    }

    if (showUnit) {
      value += ' GO';
    }
    return value;
  }
}

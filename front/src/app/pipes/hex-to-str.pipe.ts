import {PipeTransform, Pipe} from '@angular/core';

@Pipe({
  name: 'hex2str'
})

export class Hex2Str implements PipeTransform {

  transform(val: string, convert: boolean = true): string {
    if (convert) {
      let tempstr = '';
      let b = 0;
      while (b < val.length) {
        tempstr = tempstr + String.fromCharCode(parseInt(val.substr(b, 2), 16));
        b = b + 2;
      }
      return tempstr;
    }
    return val;
  }
}

import {PipeTransform, Pipe} from '@angular/core';
import {convertWithDecimals} from '../utils/functions';

@Pipe({
    name: 'weiToWholeGO'
})

export class WeiToWholeGOPipe implements PipeTransform {

    transform(val: string | number): string {
        return convertWithDecimals(val, false, false, 18).split('.')[0];
    }
}

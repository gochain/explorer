import { PipeTransform, Pipe } from '@angular/core';

@Pipe({
    name: 'weiToGO'
})

export class WeiToGOPipe implements PipeTransform {

    transform(val: string): string {        
        var moveTo = 18
        var parts = val.toString().split(".");
        if (parts[0].length > moveTo) {
            parts[0] = parts[0].slice(0, parts[0].length - moveTo) + "." + parts[0].slice(parts[0].length - moveTo, parts[0].length)
        }
        else {
            parts[0] = "0." + "0".repeat(moveTo - parts[0].length) + parts[0]
        }
        return parts.join("").toString().replace(/\.?0+$/, '') + " GO";


    }
}
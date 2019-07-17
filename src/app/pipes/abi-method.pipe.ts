import {PipeTransform, Pipe} from '@angular/core';
import {AbiItem} from 'web3-utils';

@Pipe({
  name: 'abiMethod'
})

export class AbiMethodPipe implements PipeTransform {

  transform(val: AbiItem): string {
    if (!val) {
      return null;
    }
    const inputs = val.inputs.map(input => input.name);
    // const outputs = val.outputs.map(output => output.name + (output.name ? ' ' : '') + output.type);
    return `${val.name}(${inputs.join(', ')})`;
  }
}

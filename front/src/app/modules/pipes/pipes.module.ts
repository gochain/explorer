import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BigNumberPipe} from '../../pipes/big-number.pipe';
import {WeiToGOPipe} from '../../pipes/wei-to-go.pipe';
import {Hex2Str} from '../../pipes/hex-to-str.pipe';
import {ToGweiPipe} from '../../pipes/to-gwei.pipe';
import {TrimExtra} from '../../pipes/trim-extra-data.pipe';
import {AbiMethodPipe} from '../../pipes/abi-method.pipe';

@NgModule({
  declarations: [
    BigNumberPipe,
    WeiToGOPipe,
    TrimExtra,
    ToGweiPipe,
    Hex2Str,
    AbiMethodPipe
  ],
  imports: [
    CommonModule,
  ],
  exports: [
    BigNumberPipe,
    WeiToGOPipe,
    Hex2Str,
    ToGweiPipe,
    TrimExtra,
    AbiMethodPipe,
  ]
})
export class PipesModule {
}

import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BigNumberPipe} from '../../pipes/big_number';
import {WeiToGOPipe} from '../../pipes/wei_to_go';
import {Hex2Str} from '../../pipes/hex_to_str';
import {ToGwei} from '../../pipes/to-gwei';
import {TrimExtra} from '../../pipes/trim-extra-data';

@NgModule({
  declarations: [
    BigNumberPipe,
    WeiToGOPipe,
    TrimExtra,
    ToGwei,
    Hex2Str,
  ],
  imports: [
    CommonModule,
  ],
  exports: [
    BigNumberPipe,
    WeiToGOPipe,
    Hex2Str,
    ToGwei,
    TrimExtra,
  ]
})
export class PipesModule {
}

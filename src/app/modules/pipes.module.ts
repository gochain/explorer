import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BigNumberPipe} from '../pipes/big_number';
import {WeiToGOPipe} from '../pipes/wei_to_go';
import {ToGwei} from '../pipes/to-gwei';
import {TrimExtra} from '../pipes/trim-extra-data';

@NgModule({
  declarations: [
    BigNumberPipe,
    WeiToGOPipe,
    TrimExtra,
    ToGwei,
  ],
  imports: [
    CommonModule,
  ],
  exports: [
    BigNumberPipe,
    WeiToGOPipe,
    ToGwei,
    TrimExtra,
  ]
})
export class PipesModule {
}

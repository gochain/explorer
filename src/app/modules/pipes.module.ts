import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BigNumberPipe} from '../pipes/big_number';
import {WeiToGOPipe} from '../pipes/wei_to_go';
import {ToGwei} from '../pipes/to-gwei';

@NgModule({
  declarations: [
    BigNumberPipe,
    WeiToGOPipe,
    ToGwei,
  ],
  imports: [
    CommonModule,
  ],
  exports: [
    BigNumberPipe,
    WeiToGOPipe,
    ToGwei,
  ]
})
export class PipesModule {
}

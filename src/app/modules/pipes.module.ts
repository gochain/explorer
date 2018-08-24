import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {BigNumberPipe} from '../pipes/big_number';
import {WeiToGOPipe} from '../pipes/wei_to_go';


@NgModule({
  declarations: [
    BigNumberPipe,
    WeiToGOPipe,
  ],
  imports: [
    CommonModule,
  ],
  exports: [
    BigNumberPipe,
    WeiToGOPipe,
  ]
})
export class PipesModule {
}
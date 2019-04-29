/*CORE*/
import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {RouterModule} from '@angular/router';
/*COMPONENTS*/
import {WalletUseComponent} from './wallet-use/wallet-use.component';
/*SERVICES*/
import {WalletService} from './wallet.service';
import {PipesModule} from '../pipes/pipes.module';

@NgModule({
  declarations: [
    WalletUseComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    PipesModule,
    ReactiveFormsModule,
    RouterModule,
  ],
  providers: [WalletService],
  exports: [
    WalletUseComponent,
  ],
})
export class WalletCommonModule {
}

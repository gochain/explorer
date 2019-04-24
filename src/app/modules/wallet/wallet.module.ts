/*CORE*/
import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
/*MODULES*/
import {WalletRoutingModule} from './wallet-routing.module';
import {TabsModule} from '../tabs/tabs.module';
/*COMPONENTS*/
import {WalletMainComponent} from './wallet-main/wallet-main.component';
import {WalletComponent} from './wallet/wallet.component';
import {WalletCreateComponent} from './wallet-create/wallet-create.component';
import {WalletSendComponent} from './wallet-send/wallet-send.component';
/*SERVICES*/
import {WalletService} from './wallet.service';
import {PipesModule} from '../pipes/pipes.module';
import {WalletCommonModule} from './wallet-common.module';

@NgModule({
  declarations: [
    WalletComponent,
    WalletMainComponent,
    WalletCreateComponent,
    WalletSendComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    PipesModule,
    ReactiveFormsModule,
    TabsModule,
    WalletRoutingModule,
    WalletCommonModule,
  ],
  providers: [WalletService],
})
export class WalletModule {
}

import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';

import {WalletRoutingModule} from './wallet-routing.module';
import {WalletMainComponent} from './wallet-main/wallet-main.component';
import {WalletComponent} from './wallet/wallet.component';
import {WalletCreateComponent} from './wallet-create/wallet-create.component';
import {WalletSendComponent} from './wallet-send/wallet-send.component';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {TabsModule} from '../tabs/tabs.module';
import {WalletService} from './wallet.service';

@NgModule({
  declarations: [WalletComponent, WalletMainComponent, WalletCreateComponent, WalletSendComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TabsModule,
    WalletRoutingModule,
  ],
  providers: [WalletService],
})
export class WalletModule {
}

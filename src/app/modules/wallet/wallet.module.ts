/*CORE*/
import {NgModule} from '@angular/core';
/*MODULES*/
import {WalletRoutingModule} from './wallet-routing.module';
/*COMPONENTS*/
import {WalletMainComponent} from './wallet-main/wallet-main.component';
import {WalletComponent} from './wallet/wallet.component';
import {WalletCreateComponent} from './wallet-create/wallet-create.component';
import {WalletAccountComponent} from './wallet-account/wallet-account.component';
/*SERVICES*/
import {WalletGuard} from '../../guards/wallet.guard';
import {DeployerComponent} from './deployer/deployer.component';
import {SenderComponent} from './sender/sender.component';
import {SharedModule} from '../shared/shared.module';
import {WalletSharedModule} from './wallet-shared.modules';
import {WalletService} from './wallet.service';

@NgModule({
  declarations: [
    WalletComponent,
    WalletMainComponent,
    WalletCreateComponent,
    WalletAccountComponent,
    DeployerComponent,
    SenderComponent,
  ],
  imports: [
    SharedModule,
    WalletSharedModule,
    WalletRoutingModule,
  ],
  providers: [WalletService, WalletGuard],
})
export class WalletModule {
}

import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {WalletMainComponent} from './wallet-main/wallet-main.component';
import {WalletCreateComponent} from './wallet-create/wallet-create.component';
import {WalletSendComponent} from './wallet-send/wallet-send.component';
import {WalletComponent} from './wallet/wallet.component';

const routes: Routes = [
  {
    path: '',
    component: WalletComponent,
    children: [
      {path: '', component: WalletMainComponent},
      {path: 'create', component: WalletCreateComponent},
      {path: 'send', component: WalletSendComponent},
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class WalletRoutingModule {
}

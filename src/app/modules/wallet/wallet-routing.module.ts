/*CORE*/
import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
/*SERVICES*/
import {WalletGuard} from '../../guards/wallet.guard';
/*COMPONENTS*/
import {WalletMainComponent} from './wallet-main/wallet-main.component';
import {WalletCreateComponent} from './wallet-create/wallet-create.component';
import {WalletAccountComponent} from './wallet-account/wallet-account.component';
import {WalletComponent} from './wallet/wallet.component';
import {WalletAccountComponentt} from './wallet-account1/wallet-account-componentt.component';

const routes: Routes = [
  {
    path: '',
    component: WalletComponent,
    children: [
      {path: '', component: WalletMainComponent},
      {path: 'create', component: WalletCreateComponent},
      {
        path: 'send',
        redirectTo: 'account',
      },
      {
        path: 'account',
        component: WalletAccountComponent,
        canActivate: [WalletGuard],
      },
      // {path: 'use', component: WalletAccountComponentt},
    ]
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class WalletRoutingModule {
}

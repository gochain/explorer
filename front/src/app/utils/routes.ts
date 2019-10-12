/*CORE*/
import {Routes} from '@angular/router';
/*COMPONENTS*/
import {BlockComponent} from '../scenes/block/block.component';
import {TransactionComponent} from '../scenes/transaction/transaction.component';
import {AddressComponent} from '../scenes/address/address.component';
import {RichlistComponent} from '../scenes/richlist/richlist.component';
import {HomeComponent} from '../scenes/home/home.component';
import {PageNotFoundComponent} from '../scenes/page-not-found/page-not-found.component';
import {ContractComponent} from '../scenes/contract/contract.component';
import {TokenAssetComponent} from '../scenes/token-asset/token-asset.component';
import {WalletMainComponent} from '../scenes/wallet-main/wallet-main.component';
import {WalletCreateComponent} from '../scenes/wallet-create/wallet-create.component';
import {WalletAccountComponent} from '../scenes/wallet-account/wallet-account.component';
import {SignersComponent} from '../scenes/signers/signers.component';
import {ContractsComponent} from '../scenes/contracts/contracts.component';
/*SERVICES*/
import {CommonService} from '../services/common.service';
import {WalletGuard} from '../guards/wallet.guard';
/*UTILS*/
import {ROUTES} from './constants';

export const APP_ROUTES: Routes = [
  {
    path: '',
    resolve: {rpcProvider: CommonService},
    children: [
      {
        path: 'wallet',
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
        ]
      },
      {path: ROUTES.BLOCK + '/:id', component: BlockComponent},
      {
        path: ROUTES.TRANSACTION + '/:id',
        component: TransactionComponent,
      },
      {
        path: ROUTES.ADDRESS_FULL + '/:id',
        component: AddressComponent,
      },
      {
        path: ROUTES.ADDRESS_FULL + '/:id/tx/:nonce_id',
        component: TransactionComponent,
      },

      {
        path: ROUTES.ADDRESS + '/:id',
        component: AddressComponent,
      },
      {
        path: ROUTES.ADDRESS + '/:id/tx/:nonce_id',
        component: TransactionComponent,
      },

      {
        path: ROUTES.TOKEN + '/:id',
        component: AddressComponent,
      },
      {
        path: ROUTES.TOKEN + '/:id/asset/:tokenId',
        component: TokenAssetComponent,
      },
      {path: ROUTES.VERIFY, component: ContractComponent},
      {path: ROUTES.RICHLIST, component: RichlistComponent},
      {path: ROUTES.SIGNERS, component: SignersComponent},
      {path: ROUTES.CONTRACTS, component: ContractsComponent},
      /*{path: ROUTES.SETTINGS, component: SettingsComponent},*/
      {path: ROUTES.HOME, component: HomeComponent},
      {path: '', pathMatch: 'full', redirectTo: 'home'},
    ],
  },
  {path: '**', component: PageNotFoundComponent}
];

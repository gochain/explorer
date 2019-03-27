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
import {TokenComponent} from '../scenes/token/token.component';
// import {SettingsComponent} from '../scenes/settings/settings.component';
/*UTILS*/
import {ROUTES} from './constants';
import {CommonService} from '../services/common.service';

export const APP_ROUTES: Routes = [
  {
    path: 'wallet',
    loadChildren: './modules/wallet/wallet.module#WalletModule',
    resolve: {rpcProvider: CommonService},
  },
  {path: ROUTES.BLOCK + '/:id', component: BlockComponent},
  {path: ROUTES.TRANSACTION + '/:id', component: TransactionComponent},
  {path: ROUTES.ADDRESS_FULL + '/:id', component: AddressComponent},
  {path: ROUTES.ADDRESS + '/:id', component: AddressComponent},
  {path: ROUTES.TOKEN + '/:id', component: TokenComponent},
  {
    path: ROUTES.TOKEN + '/:id/asset/:tokenId',
    component: TokenAssetComponent,
    resolve: {rpcProvider: CommonService},
  },
  {path: ROUTES.VERIFY, component: ContractComponent},
  {path: ROUTES.RICHLIST, component: RichlistComponent},
  /*{path: ROUTES.SETTINGS, component: SettingsComponent},*/
  {path: ROUTES.HOME, component: HomeComponent},
  {path: '', pathMatch: 'full', redirectTo: 'home'},
  {path: '**', component: PageNotFoundComponent}
];

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
import {WalletComponent} from '../scenes/wallet/wallet.component';
import {CreateWalletComponent} from '../scenes/create-wallet/create-wallet.component';
import {SendTxComponent} from '../scenes/send-tx/send-tx.component';
// import {SettingsComponent} from '../scenes/settings/settings.component';
/*UTILS*/
import {ROUTES} from './constants';

export const APP_ROUTES: Routes = [
  {path: ROUTES.BLOCK + '/:id', component: BlockComponent},
  {path: ROUTES.TRANSACTION + '/:id', component: TransactionComponent},
  {path: ROUTES.ADDRESS + '/:id', component: AddressComponent},
  {path: ROUTES.VERIFY, component: ContractComponent},
  {path: ROUTES.RICHLIST, component: RichlistComponent},
  {path: ROUTES.WALLET, component: WalletComponent},
  {path: ROUTES.CREATE_WALLET, component: CreateWalletComponent},
  {path: ROUTES.SEND_TX, component: SendTxComponent},
  /*{path: ROUTES.SETTINGS, component: SettingsComponent},*/
  {path: ROUTES.HOME, component: HomeComponent},
  {path: '', pathMatch: 'full', redirectTo: 'home'},
  {path: '**', component: PageNotFoundComponent}
];

import {Routes} from '@angular/router';
import {BlockComponent} from '../scenes/block/block.component';
import {TransactionComponent} from '../scenes/transaction/transaction.component';
import {AddressComponent} from '../scenes/address/address.component';
import {RichlistComponent} from '../scenes/richlist/richlist.component';
import {HomeComponent} from '../scenes/home/home.component';
import {PageNotFoundComponent} from '../scenes/page-not-found/page-not-found.component';

export const ROUTES = {
  HOME: 'home',
  BLOCK: 'block',
  ADDRESS: 'addr',
  RICHLIST: 'richlist',
  TRANSACTION: 'tx',
};

export const APP_ROUTES: Routes = [
  {path: ROUTES.BLOCK + '/:id', component: BlockComponent},
  {path: ROUTES.TRANSACTION + '/:id', component: TransactionComponent},
  {path: ROUTES.ADDRESS + '/:id', component: AddressComponent},
  {path: ROUTES.RICHLIST, component: RichlistComponent},
  {path: ROUTES.HOME, component: HomeComponent},
  {path: '', pathMatch: 'full', redirectTo: 'home'},
  {path: '**', component: PageNotFoundComponent}
];

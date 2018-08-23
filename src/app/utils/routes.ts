import {Routes} from '@angular/router';
import {BlockComponent} from '../scenes/block/block.component';
import {TransactionComponent} from '../scenes/transaction/transaction.component';
import {AddressComponent} from '../scenes/address/address.component';
import {RichlistComponent} from '../scenes/richlist/richlist.component';
import {HomeComponent} from '../scenes/home/home.component';
import {PageNotFoundComponent} from '../scenes/page-not-found/page-not-found.component';

export const APP_ROUTES: Routes = [
  {path: 'block/:id', component: BlockComponent},
  {path: 'tx/:id', component: TransactionComponent},
  {path: 'address/:id', component: AddressComponent},
  {path: 'richlist', component: RichlistComponent},
  // { path: 'send-tx', component: SendTxComponent },
  {
    path: '',
    pathMatch: 'full',
    component: HomeComponent
  },
  {path: '**', component: PageNotFoundComponent}
];
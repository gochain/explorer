/*CORE*/
import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
/*COMPONENTS*/
import {AppComponent} from './app.component';
import {BlockComponent} from './scenes/block/block.component';
import {TransactionComponent} from './scenes/transaction/transaction.component';
import {AddressComponent} from './scenes/address/address.component';
import {HomeComponent} from './scenes/home/home.component';
import {PageNotFoundComponent} from './scenes/page-not-found/page-not-found.component';
import {RichlistComponent} from './scenes/richlist/richlist.component';
import {HeaderComponent} from './components/header/header.component';
import {SearchComponent} from './components/search/search.component';
import {PaginationComponent} from './components/pagination/pagination.component';
import {TokenAssetComponent} from './scenes/token-asset/token-asset.component';
import {OwnedTokensComponent} from './components/owned-tokens/owned-tokens.component';
import {InfoComponent} from './components/info/info.component';
import {ContractComponent} from './scenes/contract/contract.component';
import {SignersComponent} from './scenes/signers/signers.component';
import {ContractsComponent} from './scenes/contracts/contracts.component';
import {AddrTransactionsComponent} from './components/addr-transactions/addr-transactions.component';
import {AddrInternalTxsComponent} from './components/addr-internal-txs/addr-internal-txs.component';
import {ContractSourceComponent} from './components/contract-source/contract-source.component';
import {TokenTxsComponent} from './components/token-txs/token-txs.component';
import {TokenHoldersComponent} from './components/token-holders/token-holders.component';

import {WalletComponent} from './scenes/wallet/wallet.component';
import {WalletMainComponent} from './scenes/wallet-main/wallet-main.component';
import {WalletCreateComponent} from './scenes/wallet-create/wallet-create.component';
import {WalletAccountComponent} from './scenes/wallet-account/wallet-account.component';
import {DeployerComponent} from './components/deployer/deployer.component';
import {SenderComponent} from './components/sender/sender.component';
import {InteractorComponent} from './components/interactor/interactor.component';
/*SERVICES*/
import {ApiService} from './services/api.service';
import {CommonService} from './services/common.service';
import {LayoutService} from './services/layout.service';
import {WalletService} from './services/wallet.service';
import {MetaService} from './services/meta.service';
/*GUARDS*/
import {WalletGuard} from './guards/wallet.guard';
/*MODULES*/
import {TabsModule} from './modules/tabs/tabs.module';
import {PipesModule} from './modules/pipes/pipes.module';
import {NgProgressModule} from '@ngx-progressbar/core';
import {NgProgressHttpModule} from '@ngx-progressbar/http';
import {ToastrModule} from './modules/toastr/toastr.module';
import {ViewportSizeModule} from './modules/viewport-size/viewport-size.module';
import {PieChartModule} from '@swimlane/ngx-charts';
import {DataTableModule} from 'ng-angular8-datatable';
import {ClipboardModule} from 'ngx-clipboard';
/*PIPES*/
import {TimeAgoPipe} from 'time-ago-pipe';
/*UTILS*/
import {APP_ROUTES} from './utils/routes';
import {APP_BASE_HREF, CommonModule} from '@angular/common';
import {VIEWPORT_SIZES} from './modules/viewport-size/contants';

@NgModule({
  declarations: [
    AppComponent,
    BlockComponent,
    TransactionComponent,
    AddressComponent,
    HomeComponent,
    PageNotFoundComponent,
    TimeAgoPipe,
    RichlistComponent,
    HeaderComponent,
    SearchComponent,
    PaginationComponent,
    InfoComponent,
    ContractComponent,
    TokenAssetComponent,
    OwnedTokensComponent,
    AddrTransactionsComponent,
    AddrInternalTxsComponent,
    ContractSourceComponent,
    TokenTxsComponent,
    TokenHoldersComponent,
    WalletComponent,
    WalletMainComponent,
    WalletCreateComponent,
    WalletAccountComponent,
    DeployerComponent,
    SenderComponent,
    InteractorComponent,
    SignersComponent,
    ContractsComponent,
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    PipesModule,
    TabsModule,
    DataTableModule,
    ViewportSizeModule.forRoot(VIEWPORT_SIZES),
    NgProgressModule.withConfig({
      trickleSpeed: 200,
      min: 20,
      meteor: false,
      spinner: false
    }),
    NgProgressHttpModule,
    ToastrModule.forRoot(),
    RouterModule.forRoot(APP_ROUTES),
    BrowserModule,
    BrowserAnimationsModule,
    PieChartModule,
    ClipboardModule,
  ],
  providers: [
    {provide: APP_BASE_HREF, useValue: '/'},
    ApiService,
    CommonService,
    LayoutService,
    MetaService,
    WalletService,
    WalletGuard,
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

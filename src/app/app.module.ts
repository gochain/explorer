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

import {AddrTransactionsComponent} from './components/addr-transactions/addr-transactions.component';
import {AddrInternalTxsComponent} from './components/addr-internal-txs/addr-internal-txs.component';
import {ContractSourceComponent} from './components/contract-source/contract-source.component';
import {TokenTxsComponent} from './components/token-txs/token-txs.component';
import {TokenHoldersComponent} from './components/token-holders/token-holders.component';

import {WalletComponent} from './modules/wallet/wallet/wallet.component';
import {WalletMainComponent} from './modules/wallet/wallet-main/wallet-main.component';
import {WalletCreateComponent} from './modules/wallet/wallet-create/wallet-create.component';
import {WalletAccountComponent} from './modules/wallet/wallet-account/wallet-account.component';
import {DeployerComponent} from './modules/wallet/deployer/deployer.component';
import {SenderComponent} from './modules/wallet/sender/sender.component';
import {InteractorComponent} from './modules/wallet/interactor/interactor.component';
/*SERVICES*/
import {ApiService} from './services/api.service';
import {CommonService} from './services/common.service';
import {LayoutService} from './services/layout.service';
import {WalletService} from './modules/wallet/wallet.service';
import {MetaService} from './services/meta.service';
/*GUARDS*/
import {WalletGuard} from './guards/wallet.guard';
/*MODULES*/
import {TabsModule} from './modules/tabs/tabs.module';
import {PipesModule} from './modules/pipes/pipes.module';
import {DirectiveModule} from './directives/directives.module';
import {NgProgressModule} from '@ngx-progressbar/core';
import {NgProgressHttpModule} from '@ngx-progressbar/http';
import {ToastrModule} from './modules/toastr/toastr.module';
import {ViewportSizeModule} from './modules/viewport-size/viewport-size.module';
import {PieChartModule} from '@swimlane/ngx-charts';
import { DataTableModule } from 'ng-angular8-datatable';
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
    /*LoaderComponent,*/
    PaginationComponent,
    /*SettingsComponent,*/
    /*ToggleSwitchComponent,*/
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
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    PipesModule,
    DirectiveModule,
    TabsModule,
    DataTableModule,
    /*SliderModule,*/
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
    BrowserModule.withServerTransition({ appId: 'serverApp' }),
    BrowserAnimationsModule,
    PieChartModule,
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

/*CORE*/
import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {RouterModule} from '@angular/router';
import {FormsModule} from '@angular/forms';
import {HttpClientModule} from '@angular/common/http';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FlexLayoutModule} from '@angular/flex-layout';
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
import {LoaderComponent} from './components/loader/loader.component';
import {PaginationComponent} from './components/pagination/pagination.component';
import { MobileMenuComponent } from './components/mobile-menu/mobile-menu.component';
import { MobileSearchComponent } from './mobile-search/mobile-search.component';
import { ToggleSwitchComponent } from './components/toggle-switch/toggle-switch.component';
import { MobileHeaderComponent } from './components/mobile-header/mobile-header.component';
// import {SettingsComponent} from './scenes/settings/settings.component';
/*SERVICES*/
import {ApiService} from './services/api.service';
import {CommonService} from './services/common.service';
import {LayoutService} from './services/layout.service';
import {ViewportSizeModule} from './modules/viewport-size/viewport-size.module';
/*MODULES*/
import {TabsModule} from './modules/tabs/tabs.module';
import {PipesModule} from './modules/pipes.module';
import {DirectiveModule} from './directives/directives.module';
/*PIPES*/
import {TimeAgoPipe} from 'time-ago-pipe';
/*UTILS*/
import {APP_ROUTES} from './utils/routes';
import {APP_BASE_HREF} from '@angular/common';
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
    LoaderComponent,
    PaginationComponent,
    /*SettingsComponent,*/
    ToggleSwitchComponent,
    MobileHeaderComponent,
    MobileMenuComponent,
    MobileSearchComponent,
  ],
  imports: [
    RouterModule.forRoot(APP_ROUTES),
    BrowserModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    FormsModule,
    HttpClientModule,
    PipesModule,
    DirectiveModule,
    ViewportSizeModule.forRoot(VIEWPORT_SIZES),
    TabsModule
  ],
  providers: [
    {provide: APP_BASE_HREF, useValue: '/'},
    ApiService,
    CommonService,
    LayoutService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}

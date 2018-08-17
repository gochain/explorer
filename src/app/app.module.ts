import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpClientModule, HttpClient } from '@angular/common/http';

import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FlexLayoutModule } from '@angular/flex-layout';
import {
  MatFormFieldModule, MatButtonModule, MatCheckboxModule, MatInputModule, MatCardModule, MatProgressSpinnerModule, MatProgressBarModule, MatToolbarModule, MatSnackBarModule,
  MatSelectModule, MatListModule, MatIconModule, MatSidenavModule, MatMenuModule
} from '@angular/material';
import { TimeAgoPipe } from 'time-ago-pipe';
import {BigNumberPipe} from './big_number'
import {WeiToGOPipe} from './wei_to_go'

import { AppComponent } from './app.component';

import { BlockComponent } from './block/block.component';
import { TransactionComponent } from './transaction/transaction.component';
import { AddressComponent } from './address/address.component';
import { HomeComponent } from './home/home.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';
import { RichlistComponent } from './richlist/richlist.component';

const appRoutes: Routes = [
  { path: 'block/:id', component: BlockComponent },
  { path: 'tx/:id', component: TransactionComponent },
  { path: 'address/:id', component: AddressComponent },
  { path: 'richlist', component: RichlistComponent },
  // { path: 'send-tx', component: SendTxComponent },
  // { path: 'hero/:id',      component: HeroDetailComponent },
  // {
  //   path: 'heroes',
  //   component: HeroListComponent,
  //   data: { title: 'Heroes List' }
  // },
  {
    path: '',
    // redirectTo: '/heroes',
    pathMatch: 'full',
    component: HomeComponent
  },
  { path: '**', component: PageNotFoundComponent }
];

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
    BigNumberPipe,
    WeiToGOPipe
  ],
  imports: [
    RouterModule.forRoot(
      appRoutes,
      { enableTracing: true } // <-- debugging purposes only
    ),
    BrowserModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    FormsModule,
    MatFormFieldModule, MatButtonModule, MatCheckboxModule, MatInputModule, MatCardModule, MatProgressSpinnerModule, MatProgressBarModule, MatToolbarModule, MatSnackBarModule,
    MatSelectModule, MatListModule, MatIconModule, MatSidenavModule, MatMenuModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }

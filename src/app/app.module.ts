import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { HttpClientModule, HttpClient } from '@angular/common/http';

import { AngularFireModule } from 'angularfire2';
import { AngularFirestoreModule } from 'angularfire2/firestore';
import { AngularFireAuthModule } from 'angularfire2/auth';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {FlexLayoutModule} from '@angular/flex-layout';
import {MatFormFieldModule, MatButtonModule, MatCheckboxModule, MatInputModule, MatCardModule, MatProgressSpinnerModule, MatProgressBarModule, MatToolbarModule, MatSnackBarModule,
        MatSelectModule, MatListModule, MatIconModule} from '@angular/material';
import {TimeAgoPipe} from 'time-ago-pipe';

import { AppComponent } from './app.component';

import { environment } from '../environments/environment';
import { BlockComponent } from './block/block.component';
import { TransactionComponent } from './transaction/transaction.component';
import { AddressComponent } from './address/address.component';
import { HomeComponent } from './home/home.component';
import { PageNotFoundComponent } from './page-not-found/page-not-found.component';

const appRoutes: Routes = [
  { path: 'block/:id', component: BlockComponent },
  { path: 'tx/:id', component: TransactionComponent },
  { path: 'address/:id', component: AddressComponent },
  // { path: 'send-tx', component: SendTxComponent },
  // { path: 'hero/:id',      component: HeroDetailComponent },
  // {
  //   path: 'heroes',
  //   component: HeroListComponent,
  //   data: { title: 'Heroes List' }
  // },
  { path: '',
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
    TimeAgoPipe
  ],
  imports: [
    RouterModule.forRoot(
      appRoutes,
      { enableTracing: true } // <-- debugging purposes only
    ),
    BrowserModule,
    AngularFireModule.initializeApp(environment.firebase),
    AngularFirestoreModule, // imports firebase/firestore, only needed for database features
    AngularFireAuthModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    MatFormFieldModule, MatButtonModule, MatCheckboxModule, MatInputModule, MatCardModule, MatProgressSpinnerModule, MatProgressBarModule, MatToolbarModule, MatSnackBarModule,
    MatSelectModule, MatListModule, MatIconModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }

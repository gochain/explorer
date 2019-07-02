/*MODULES*/
import {NgModule} from '@angular/core';
/*COMPONENTS*/
import {ContractInteractorComponent} from '../wallet/contract-interactor/contract-interactor.component';
import {SharedModule} from '../shared/shared.module';
import {WalletService} from './wallet.service';

@NgModule({
  declarations: [
    ContractInteractorComponent,
  ],
  imports: [
    SharedModule
  ],
  providers: [
    WalletService
  ],
  exports: [
    ContractInteractorComponent,
  ],
})
export class WalletSharedModule {
}

/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {Subscription} from 'rxjs';
import {flatMap} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {MetaService} from '../../services/meta.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Address} from '../../models/address.model';
/*UTILS*/
import {META_TITLES} from '../../utils/constants';
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';

@Component({
  selector: 'app-wallet-account',
  templateUrl: './wallet-account.component.html',
  styleUrls: ['./wallet-account.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class WalletAccountComponent implements OnInit, OnDestroy {
  accountAddr: Address;
  private _subsArr$: Subscription[] = [];

  constructor(
    public walletService: WalletService,
    private _metaService: MetaService,
    private _commonService: CommonService,
    private _router: Router,
  ) {
  }

  ngOnInit(): void {
    this._metaService.setTitle(META_TITLES.WALLET.title);
    this._subsArr$.push(this.walletService.logged$.pipe(
      flatMap(() => this._commonService.getAddress(this.walletService.accountAddress)),
    ).subscribe((addr: Address) => {
      this.accountAddr = addr;
    }));
  }

  closeWallet(): void {
    // wallet service close account will be called in ngOnDestroy
    this._router.navigate(['wallet']);
  }


  ngOnDestroy(): void {
    this.walletService.resetProcessing();
    this.walletService.closeAccount();
  }
}

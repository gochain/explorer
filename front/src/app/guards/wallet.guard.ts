/*CORE*/
import {Injectable} from '@angular/core';
import {CanActivate, Router} from '@angular/router';
/*SERVICES*/
import {WalletService} from '../services/wallet.service';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class WalletGuard implements CanActivate {
  constructor(
    private _router: Router,
    private _walletService: WalletService,
  ) {
  }

  canActivate(): Observable<boolean> {
    return this._walletService.logged$;
    /*if (!this._walletService.account) {
      this._router.navigate(['wallet']);
      return false;
    }
    return true;*/
  }
}

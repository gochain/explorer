/*CORE*/
import {Injectable} from '@angular/core';
import {CanActivate, Router} from '@angular/router';
import {Observable} from 'rxjs';
import {filter, mergeMap} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../services/wallet.service';

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
    return this._walletService.ready$.pipe(
      mergeMap(() => this._walletService.logged$),
      filter((logged: boolean) => {
        if (!logged) {
          this._router.navigate(['wallet']);
        }
        return logged;
      }),
    );
  }
}

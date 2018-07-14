import { Injectable } from '@angular/core';
import { environment } from '../environments/environment';

@Injectable()
export class Globals {
  network: string = 'mainnet';

  public rpcHost(): string {
    if (this.network == "testnet") {
      return "https://testnet-rpc.gochain.io";
    }
    return "https://rpc.gochain.io";
  }

  public explorerHost(): string {
    if (environment.production) {
      if (this.network == "testnet") {
        return "https://testnet-explorer.gochain.io";
      }
      return "https://explorer.gochain.io";
    }
    return 'http://localhost:8000';
  }

  public apiHost(): string {
    if (environment.production) {
      if (this.network == "testnet") {
        return "https://testnet-explorer.gochain.io";
      }
      return "https://explorer.gochain.io";
    }
    return 'http://localhost:8000';
  }
}

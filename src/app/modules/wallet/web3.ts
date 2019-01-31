import {InjectionToken} from '@angular/core';
import Web3 from 'web3';

const HTTP_PROVIDER = /^explorer\.gochain\.io/.test(location.hostname) ? 'https://rpc.gochain.io' : 'https://testnet-rpc.gochain.io';

export const WEB3 = new InjectionToken<Web3>('web3', {
  providedIn: 'root',
  factory: () => {
    try {
      if (Web3) {
        const provider = new Web3.providers.HttpProvider(HTTP_PROVIDER);
        return new Web3(provider);
      } else {
        console.log('No web3? You should consider trying MetaMask!');
        return null;
      }
    } catch (err) {
      throw new Error(err);
    }
  }
});

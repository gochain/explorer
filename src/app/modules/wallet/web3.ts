import {InjectionToken} from '@angular/core';
import Web3 from 'web3';

export const WEB3 = new InjectionToken<Web3>('web3', {
  providedIn: 'root',
  factory: () => {
    try {
      if (Web3) {
        const provider = new Web3.providers.HttpProvider('https://testnet-rpc.gochain.io');
        return new Web3(provider);
      } else {
        console.log('No web3? You should consider trying MetaMask!');
      }
    } catch (err) {
      throw new Error(err);
    }
  }
});

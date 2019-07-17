/*
import {InjectionToken} from '@angular/core';
import Web3 from 'web3';

export const WEB3 = new InjectionToken<Web3>('web3', {
  providedIn: 'root',
  factory: () => {
    try {
      if (Web3) {
        return new Web3();
      } else {
        console.log('No web3? You should consider trying MetaMask!');
        return null;
      }
    } catch (err) {
      throw new Error(err);
    }
  }
});
*/

/*CORE*/
import {Subscription} from 'rxjs';
/*MODELS*/
import {ABIDefinition} from 'web3/eth/abi';
import {Address} from '../models/address.model';
import {Contract} from '../models/contract.model';
import {Badge} from '../models/badge.model';
/*UTILS*/
import {InterfaceName, StatusColor} from './enums';
import {TOKEN_TYPES} from './constants';
import {ContractAbi} from './types';

declare const window: any;

/**
 * clears array from subscriptions
 * @param {Array<Subscription>} arr
 */
export function unsubArr(arr: Array<Subscription>) {
  arr.forEach(item => item.unsubscribe());
}

/**
 * recursively checks object keys
 * @param obj
 * @param keys
 */
export function objHas(obj: any, keys: string): boolean {
  return !!keys.split('.').reduce((acc, key) => acc.hasOwnProperty(key) ? (acc[key] || 1) : false, obj);
}

/**
 *
 */
export function isPrivateMode(): boolean {
  const testLocalStorage = () => {
    try {
      if (localStorage.length) {
        return true;
      } else {
        localStorage.x = 1;
        localStorage.removeItem('x');
        return true;
      }
    } catch (e) {
      // Safari only enables cookie in private mode
      // if cookie is disabled then all client side storage is disabled
      // if all client side storage is disabled, then there is no point
      // in using private mode
      return !navigator.cookieEnabled;
    }
  };
  // Chrome & Opera
  if (window.webkitRequestFileSystem) {
    return window.webkitRequestFileSystem(0, 0, true, false);
  }
  // Firefox
  if ('MozAppearance' in document.documentElement.style) {
    const db = indexedDB.open('test');
    db.onerror = () => {
      return true;
    };
    db.onsuccess = () => {
      return false;
    };
  }
  // Safari
  if (/constructor/i.test(window.HTMLElement)) {
    return testLocalStorage();
  }
  // IE10+ & Edge
  if (!window.indexedDB && (window.PointerEvent || window.MSPointerEvent)) {
    return true;
  }
  // others
  return false;
}

/**
 * returns only abi methods
 * @param abi
 */
export function getAbiMethods(abi: ABIDefinition[]): ABIDefinition[] {
  return abi.filter((abiDef: ABIDefinition) => abiDef.type === 'function');
}

/**
 * forms badges for contract
 * @param address
 * @param contract
 */
export function makeContractBadges(address: Address, contract: Contract): Badge[] {
  const badges: Badge[] = [];
  if (contract.valid) {
    badges.push({
      type: StatusColor.Success,
      text: 'Verified',
    });
  }
  if (contract.abi && contract.abi.length) {
    badges.push({
      type: StatusColor.Info,
      text: 'Has ABI',
    });
  }
  address.erc_types.forEach((value: string) => {
    badges.push({
      type: StatusColor.Info,
      text: TOKEN_TYPES[value],
    });
  });

  return badges;
}

/**
 * makes contract abi
 * @param interfaceNames
 * @param abi
 */
export function makeContractAbi(interfaceNames: InterfaceName[], abi: ContractAbi): ABIDefinition[] {
  const contractAbi: ABIDefinition[] = [];
  interfaceNames.forEach((name: InterfaceName) => {
    if (abi[name]) {
      contractAbi.push(abi[name]);
    }
  });
  return contractAbi;
}

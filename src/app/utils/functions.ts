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
 * checks obj
 * @param obj
 */
export function objIsEmpty(obj: any): boolean {
  return Object.entries(obj).length === 0 && obj.constructor === Object;
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

/**
 * get appropriate data from function result
 * @param decoded
 */
export function getDecodedData(decoded: object): any[][] {
  const arrR: any[][] = [];
  // let mapR: Map<any,any> = new Map<any,any>();
  // for (let j = 0; j < decoded.__length__; j++){
  //   mapR.push([decoded[0], decoded[1]])
  // }
  Object.keys(decoded).forEach((key) => {
    // mapR[key] = decoded[key];
    if (key.startsWith('__')) {
      return;
    }
    if (!decoded[key].payable || decoded[key].constant) {
      arrR.push([key, decoded[key]]);
    }
  });
  return arrR;
}

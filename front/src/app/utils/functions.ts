/*CORE*/
import {Subscription} from 'rxjs';
/*MODELS*/
import {AbiItem} from 'web3-utils';
import {Address} from '../models/address.model';
import {Contract} from '../models/contract.model';
import {Badge} from '../models/badge.model';
/*UTILS*/
import {InterfaceName, StatusColor} from './enums';
import {TOKEN_ABI_NAMES, TOKEN_TYPES} from './constants';
import {ContractAbi} from './types';

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
export function getAbiMethods(abi: AbiItem[]): AbiItem[] {
  return abi.filter((abiDef: AbiItem) => abiDef.type === 'function');
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
export function makeContractAbi(interfaceNames: InterfaceName[], abi: ContractAbi): AbiItem[] {
  const contractAbi: AbiItem[] = [];
  interfaceNames.forEach((name: InterfaceName) => {
    if (abi[name]) {
      contractAbi.push(abi[name]);
    }
  });
  return contractAbi;
}

/**
 *
 * @param val
 * @param showUnit
 * @param removeTrailingZeros
 * @param decimals
 * @param unitName
 */
export function convertWithDecimals(
  val: string | number,
  showUnit: boolean = true,
  removeTrailingZeros: boolean = false,
  decimals: number = 18,
  unitName: string = 'GO',
): string {
  if (!val) {
    return;
  }
  const parts = val.toString().split('.');
  if (parts[0].length > decimals) {
    parts[0] =
      parts[0].slice(0, parts[0].length - decimals)
      + '.'
      + parts[0].slice(parts[0].length - decimals, parts[0].length);
  } else {
    parts[0] = '0.' + '0'.repeat(decimals - parts[0].length) + parts[0];
  }
  let value: string = parts.join('').toString();

  if (removeTrailingZeros) {
    // replace trailing zeros with exact amount of spaces
    value = value.replace(/0(?=(0+$|$))/g, ` `);
    value = value.replace(/\.(?=\s)/g, ` `);
  } else {
    // delete trailing zeros
    value = value.replace(/\.?0+$/, '');
  }

  if (showUnit) {
    value += ' ' + unitName;
  }
  // remove dot in the end
  value = value.replace(/\.$/, '');
  return value;
}

export function numberWithCommas(val: string): string {
  if (val == null) {
    return val;
  }
  const parts = val.toString().split('.');
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',');
  return parts.join('.');
}

/**
 * get appropriate data from function result
 * @param decoded
 * @param abi
 * @param addr
 */
export function getDecodedData(decoded: object, abi: AbiItem, addr: Address): any[][] {
  const arrR: any[][] = [];
  // let mapR: Map<any,any> = new Map<any,any>();
  // for (let j = 0; j < decoded.__length__; j++){
  //   mapR.push([decoded[0], decoded[1]])
  // }
  Object.keys(decoded).forEach((key) => {
    let val = decoded[key];
    if (addr && addr.decimals && TOKEN_ABI_NAMES.includes(abi.name)) {
      val = numberWithCommas(convertWithDecimals(val, true, true, addr.decimals, addr.token_symbol));
    }
    // mapR[key] = decoded[key];
    if (key.startsWith('__')) {
      return;
    }
    if (!decoded[key].payable || decoded[key].constant) {
      arrR.push([key, val]);
    }
  });
  return arrR;
}

/**
 *
 * @param value
 */
export function isHex(val: string): boolean {
  return /^[0-9A-F]+$/i.test(val);
}

/**
 *
 * @param arr
 * @param key
 * @param desc
 */
export function sortObjArrByKey(arr: any[], key: string, desc: boolean = true) {
  if (desc) {
    arr.sort((a, b) => {
      if (a[key] > b[key]) return -1;
      if (a[key] < b[key]) return 1;
      return 0;
    });
  } else {
    arr.sort((a, b) => {
      if (a[key] > b[key]) return 1;
      if (a[key] < b[key]) return -1;
      return 0;
    });
  }
}


export function removeEmpty(obj: any): any {
  Object.keys(obj).forEach((key) => (obj[key] === null) && delete obj[key]);
}

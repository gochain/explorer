import {InterfaceName} from './enums';
import {AbiItem} from 'web3-utils';

export type ContractAbi = {
  [key in InterfaceName]: AbiItem;
};

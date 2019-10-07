import {FunctionName} from './enums';
import {AbiItem} from 'web3-utils';

export type ContractAbi = {
  [key in FunctionName]: AbiItem;
};

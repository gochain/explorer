import {FunctionName, EventID, ErcName} from './enums';
import {AbiItem} from 'web3-utils';

export type ContractAbi = {
  [key in FunctionName]: AbiItem;
};

export type ContractEventsAbi = {
  [key in EventID]: {
    [key in ErcName]: AbiItem;
  };
};

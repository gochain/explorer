import {FunctionName, EventID, ErcName} from './enums';
import {AbiItem} from 'web3-utils';

export type ContractAbi = {
  [key in FunctionName]: AbiItem;
};

export type ContractAbiByID = {
  [funcKey in string]:  AbiItem;
};

export interface AbiItemIDed extends AbiItem {
  id: string;
};

export type ContractEventsAbi = {
  [eventKey in EventID]: {
    [ercKey in ErcName]: AbiItem;
  };
};
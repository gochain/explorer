import {InterfaceName} from './enums';
import {ABIDefinition} from 'web3/eth/abi';

export type ContractAbi = {
  [key in InterfaceName]: ABIDefinition;
};

import {FunctionName} from '../utils/enums';
import {SignerDetails} from './signer-node';
import {Contract} from './contract.model';

export class Address {
  address: string;
  balance: number;
  balance_wei: number;
  decimals: number;
  token_name: string;
  token_symbol: string;
  total_supply: number;
  contract: boolean;
  erc_types: string[];
  interfaces: FunctionName[];
  ercObj: object;
  supplyOwned?: string;
  number_of_transactions: number;
  number_of_token_holders: number;
  number_of_internal_transactions: number;
  number_of_token_transactions: number;
  updated_at: Date;
  signerDetails?: SignerDetails;
  attached_contract: Contract;
  target: string;
  owner: string;
}

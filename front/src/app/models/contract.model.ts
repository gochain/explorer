import {AbiItem} from 'web3-utils';

export class Contract {
  address: string;
  byte_code: string;
  valid: boolean;
  contract_name: string;
  compiler_version: string;
  evm_version: string;
  optimization: boolean;
  source_code: string;
  abi: AbiItem[];
  created_at: Date;
  updated_at: Date;
}

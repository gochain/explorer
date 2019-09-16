import {Address} from './address.model';

export class InternalTransaction {
  contract_address: string;
  from_address: string;
  to_address: string;
  value: number;
  block_number: number;
  transaction_hash: string;
  updated_at: Date;
  created_at: Date;
  address: Address;
}

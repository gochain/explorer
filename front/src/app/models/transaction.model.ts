import {Address} from './address.model';
import {Log} from 'web3-core';

export class ProcessedLogItem {
  link?: string;
  name?: string;
  value: string;
}

export class ProcessedLogData {
  title: string;
  items: ProcessedLogItem[];
}

export class ProcessedLog {
  index: number;
  contract_address: string;
  data: ProcessedLogData[];
  removed: boolean;
  recognized: boolean;
}

export interface TxLog extends Log {
  removed: boolean;
}

export class Transaction {
  tx_hash: string;
  created_at: Date;
  value: string;
  gas_price: string;
  gas_fee: string;
  gas_limit: string;
  block_number: number;
  nonce: number;
  input_data: string;
  logs: string;
  prettifiedLogs: string;
  parsedLogs?: TxLog[];
  from: string;
  to: string;
  contract_address: string;
  status: boolean;
  addressData?: Address;
  processedLogs?: ProcessedLog[];
}

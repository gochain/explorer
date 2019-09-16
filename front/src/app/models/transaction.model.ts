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
  from: string;
  to: string;
  contract_address: string;
  status: boolean;
}

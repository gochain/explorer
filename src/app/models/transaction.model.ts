export class Transaction {
  tx_hash: string;
  created_at: Date;
  value: string;
  gas_price: string;
  gas_fee: string;
  gas_limit: string;
  block_number: number;
  nonce: string;
  input_data: string;
  to: string;
  from: string;
}
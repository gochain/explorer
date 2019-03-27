export class Address {
  address: string;
  balance: number;
  balance_wei: number;
  decimals: number;
  token_name: string;
  token_symbol: string;
  total_supply: number;
  contract: boolean;
  go20: boolean;
  erc_types: string[];
  supplyOwned?: string;
  number_of_transactions: number;
  number_of_token_holders: number;
  number_of_internal_transactions: number;
  updated_at: Date;
}

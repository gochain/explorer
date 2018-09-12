export class Address {
  address: string;
  balance: number;
  balance_wei: number;
  token_name: string;
  token_symbol: string;
  contract: boolean;
  go20: boolean;
  supplyOwned?: string;
  number_of_transactions: number;
  number_of_token_holders: number;
  number_of_internal_transactions: number;
}

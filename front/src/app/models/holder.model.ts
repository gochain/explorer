import { Decimal } from 'decimal.js/decimal';

export class Holder {
  contract_address: string;
  token_holder_address: string;
  balance: string;
  token_name: string;
  token_symbol: string;
  token_decimals: number;

  public balanceDec(): Decimal {
    var b = new Decimal(this.balance);    
    var mby = new Decimal(10).toPower(this.token_decimals);    
    return b.div(mby);
  }
}

import {ChartItem} from './chart';
import {SignerDetails} from './signer-node';

export enum SignerStatRange {
  DAILY = 'daily',
  WEEKLY = 'weekly',
  MONTHLY = 'monthly',
}

interface BlockRange {
  start_block: number;
  end_block: number;
}

export interface SignerData {
  data?: SignerDetails;
  signer_address: string;
  blocks_count: number;
  percent?: number | string;
}

export class SignerStat {
  block_range: BlockRange;
  signer_stats: SignerData[];
  range: SignerStatRange;
  totalBlocks: number;
  chartData: ChartItem[];
}

import {ChartItem} from './chart';

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
  name: string;
  url: string;
  region: string;
  signer: string;
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

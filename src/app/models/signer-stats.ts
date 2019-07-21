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

interface SignerNode {
  name: string;
  url: string;
  region: string;
}

export interface SignerData {
  signer: SignerNode;
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

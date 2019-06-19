export class Block {
  number: number;
  created_at: Date;
  hash: string;
  tx_count: number;
  parent_hash: string;
  gas_used: any;
  miner: any;
  difficulty: any;
  sha3_uncles: any;
  extra_data: any;
  nonce: number;
  gas_limit: number;
  extra_auth: boolean;
	extra_vanity: string;
	extra_has_vote: boolean;
	extra_candidate: string;
	extra_is_voter_election: boolean;
}

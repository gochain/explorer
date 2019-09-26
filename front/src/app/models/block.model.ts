import {SignerDetails} from './signer-node';

export class Extra {
  auth: boolean;
  vanity: string;
  has_vote: boolean;
  candidate: string;
  is_voter_election: boolean;
  signerDetails?: SignerDetails;
}

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
  extra: Extra;
  signerDetails?: SignerDetails;
}

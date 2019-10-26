import {SignerDetails} from './signer-node';

export class Extra {
  vanity: string;

  has_vote: boolean; // If true, then the remaining vote fields are included.
  auth: boolean; // Whether voting to authorize (add) or deauthorize (remove).
  is_voter_election: boolean; // Whether voting on voter (true) or signer (false) auth.
  candidate: string; // Candidate address.
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
  extra_data: any;
  gas_limit: number;
  extra: Extra;
  signerDetails?: SignerDetails;
}

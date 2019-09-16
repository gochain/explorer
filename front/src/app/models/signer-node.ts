export interface SignerDetails {
  name: string;
  url: string;
  region: string;
}


export class SignerNode {
  [key: string]: SignerDetails;
}


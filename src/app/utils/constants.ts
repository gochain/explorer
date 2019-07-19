import {MenuItem} from '../models/menu_item.model';
import {ErcName, InterfaceName, ThemeColor} from './enums';

export const THEME_SETTINGS = {
  color: ThemeColor.LIGHT,
};

export const LOGO_NAMES = {
  [ThemeColor.LIGHT]: 'logo_fullcolor.svg',
  [ThemeColor.DARK]: 'logo_allwhite.svg',
};

export const ROUTES = {
  HOME: 'home',
  BLOCK: 'block',
  ADDRESS_FULL: 'address',
  ADDRESS: 'addr',
  TOKEN: 'token',
  RICHLIST: 'richlist',
  TRANSACTION: 'tx',
  SETTINGS: 'settings',
  VERIFY: 'verify',
  WALLET: 'wallet',
  CREATE_WALLET: 'create-account',
  SEND_TX: 'send-tx',
  SIGNERS: 'signers',
};

export const MENU_ITEMS: MenuItem[] = [
  {
    title: 'Main',
    link: ROUTES.HOME,
  },
  {
    title: 'Rich List',
    link: ROUTES.RICHLIST,
  },
  /*{
    title: 'Verify Contract',
    link: '/verify',
  },*/
  {
    title: 'Wallet',
    link: ROUTES.WALLET,
  },
  {
    title: 'Signers',
    link: ROUTES.SIGNERS,
  },
  /*{
    title: 'Network Stats',
    link: 'https://stats.gochain.io',
    external: true
  },*/
  /*{
    title: 'Settings',
    link: '/settings',
  },*/
];

export const DEFAULT_GAS_LIMIT = 21000;

export const TOKEN_TYPES = {
  Go20: 'GO20',
  Go20Burnable: 'GO20 Burnable',
  Go20Capped: 'GO20 Capped',
  Go20Detailed: 'GO20 Detailed',
  Go20Mintable: 'GO20 Mintable',
  Go20Pausable: 'GO20 Pausable',
  Go165: 'GO165',
  Go721: 'GO721',
  Go721Burnable: 'GO721 Burnable',
  Go721Receiver: 'GO721 Receiver',
  Go721Metadata: 'GO721 Metadata',
  Go721Enumerable: 'GO721 Enumerable',
  Go721Pausable: 'GO721 Pausable',
  Go721Mintable: 'GO721 Mintable',
  Go721MetadataMintable: 'GO721 Metadata Mintable',
  Go721Full: 'GO721 Full',
  Go820: 'GO820',
  Go1155: 'GO1155',
  Go1155Receiver: 'GO1155 Receiver',
  Go1155Metadata: 'GO1155 Metadata',
  Go223: 'GO223',
  Go223Receiver: 'GO223 Receiver',
  Go621: 'GO621',
  Go777: 'GO777',
  Go777Receiver: 'GO777 Receiver',
  Go777Sender: 'GO777 Sender',
  Go827: 'GO827',
  Go884: 'GO884',
};

export const ERC_INTERFACE_IDENTIFIERS = {
  [ErcName.Go20]: [InterfaceName.Allowance, InterfaceName.Approve, InterfaceName.BalanceOf, InterfaceName.TotalSupply, InterfaceName.Transfer, InterfaceName.TransferFrom],
  [ErcName.Go721]: [InterfaceName.Approve, InterfaceName.BalanceOf, InterfaceName.GetApproved, InterfaceName.IsApprovedForAll, InterfaceName.OwnerOf, InterfaceName.SafeTransferFrom, InterfaceName.SafeTransferFrom1, InterfaceName.SetApprovalForAll, InterfaceName.TransferFrom],
};

export const TOKEN_ABI_NAMES: string[] = ['totalSupply', 'balanceOf'];

export const META_TITLES = {
  DEFAULT: {
    title: 'GoChain Explorer',
  },
  HOME: {
    title: 'Home',
  },
  BLOCK: {
    title: 'Block',
  },
  ADDRESS: {
    title: 'Address',
  },
  CONTRACT: {
    title: 'Contract',
  },
  TOKEN: {
    title: 'Token',
  },
  RICHLISLT: {
    title: 'Richlist',
  },
  TRANSACTION: {
    title: 'Transaction',
  },
  VERIFY: {
    title: 'Verify contract',
  },
  SEND_TX: {
    title: 'Send transaction',
  },
  WALLET: {
    title: 'Wallet',
  },
  CREATE_WALLET: {
    title: 'Create account',
  },
  SEND_WALLET: {
    title: 'Send GO',
  },
  DEPLOY_CONTRACT: {
    title: 'Deploy contract',
  },
  USE_CONTRACT: {
    title: 'Interact with a Smart Contract',
  },
  OPEN_WALLET: {
    title: 'Open wallet',
  },
  NOT_FOUND: {
    title: 'Not found',
  },
};

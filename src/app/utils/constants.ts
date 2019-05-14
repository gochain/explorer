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
};

export const MENU_ITEMS: MenuItem[] = [
  {
    title: 'Blocks',
    link: ROUTES.HOME,
    icon: 'fa fa-link fa-fw'
  },
  {
    title: 'Rich List',
    link: ROUTES.RICHLIST,
    icon: 'fa fa-list-ul fa-fw'
  },
  /*{
    title: 'Verify Contract',
    link: '/verify',
    icon: 'fa fa-check-square fa-fw'
  },*/
  {
    title: 'Wallet',
    link: ROUTES.WALLET,
    icon: 'fa fa-wallet fa-fw',
  },
  {
    title: 'Network Stats',
    link: 'https://stats.gochain.io',
    icon: 'fa fa-broadcast-tower fa-fw',
    external: true
  },
  /*{
    title: 'Settings',
    link: '/settings',
    icon: 'fa fa-cogs fa-fw',
  },*/
];

export const DEFAULT_GAS_LIMIT = 21000;

export const TOKEN_TYPES = {
  Erc20: 'GO20',
  Erc20Burnable: 'GO20 Burnable',
  Erc20Capped: 'GO20 Capped',
  Erc20Detailed: 'GO20 Detailed',
  Erc20Mintable: 'GO20 Mintable',
  Erc20Pausable: 'GO20 Pausable',
  Erc165: 'GO165',
  Erc721: 'GO721',
  Erc721Receiver: 'GO721 Receiver',
  Erc721Metadata: 'GO721 Metadata',
  Erc721Enumerable: 'GO721 Enumerable',
  Erc820: 'GO820',
  Erc1155: 'GO1155',
  Erc1155Receiver: 'GO1155 Receiver',
  Erc1155Metadata: 'GO1155 Metadata',
  Erc223: 'GO223',
  Erc621: 'GO621',
  Erc777: 'GO777',
  Erc777Receiver: 'GO777 Receiver',
  Erc777Sender: 'GO777 Sender',
  Erc827: 'GO827',
  Erc884: 'GO884',
};

export const ERC_INTERFACE_IDENTIFIERS = {
  [ErcName.Erc20]: [InterfaceName.Allowance, InterfaceName.Approve, InterfaceName.BalanceOf, InterfaceName.TotalSupply, InterfaceName.Transfer, InterfaceName.TransferFrom],
  [ErcName.Erc721]: [InterfaceName.Approve, InterfaceName.BalanceOf, InterfaceName.GetApproved, InterfaceName.IsApprovedForAll, InterfaceName.OwnerOf, InterfaceName.SafeTransferFrom, InterfaceName.SafeTransferFrom1, InterfaceName.SetApprovalForAll, InterfaceName.TransferFrom],
};

export const TOKEN_ABI_NAMES: string[] = ['totalSupply', 'balanceOf'];

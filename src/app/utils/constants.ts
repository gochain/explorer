import {MenuItem} from '../models/menu_item.model';
import {ThemeColor} from './enums';

export const THEME_SETTINGS = {
  color: ThemeColor.DARK,
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
  Erc20: 'Erc20',
  Erc20Burnable: 'Erc20 Burnable',
  Erc20Capped: 'Erc20 Capped',
  Erc20Detailed: 'Erc20 Detailed',
  Erc20Mintable: 'Erc20 Mintable',
  Erc20Pausable: 'Erc20 Pausable',
  Erc165: 'Erc165',
  Erc721: 'Erc721',
  Erc721Receiver: 'Erc721 Receiver',
  Erc721Metadata: 'Erc721 Metadata',
  Erc721Enumerable: 'Erc721 Enumerable',
  Erc820: 'Erc820',
  Erc1155: 'Erc1155',
  Erc1155Receiver: 'Erc1155 Receiver',
  Erc1155Metadata: 'Erc1155 Metadata',
  Erc223: 'Erc223',
  Erc621: 'Erc621',
  Erc777: 'Erc777',
  Erc777Receiver: 'Erc777 Receiver',
  Erc777Sender: 'Erc777 Sender',
  Erc827: 'Erc827',
  Erc884: 'Erc884',
};

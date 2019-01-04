import {MenuItem} from '../models/menu_item.model';

export const THEME_SETTINGS = {
  color: 'dark',
};

export const ROUTES = {
  HOME: 'home',
  BLOCK: 'block',
  ADDRESS: 'addr',
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

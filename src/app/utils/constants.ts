import {MenuItem} from '../models/menu_item.model';

export const THEME_SETTINGS = {
  color: 'dark',
};

export const MENU_ITEMS: MenuItem[] = [
  {
    title: 'Blocks',
    link: '/home',
    icon: 'fa fa-link fa-fw'
  },
  {
    title: 'Rich List',
    link: '/richlist',
    icon: 'fa fa-list-ul fa-fw'
  },
  /*{
    title: 'Verify Contract',
    link: '/verify',
    icon: 'fa fa-check-square fa-fw'
  },*/
  {
    title: 'Wallet',
    link: 'https://wallet.gochain.io',
    icon: 'fa fa-wallet fa-fw',
    external: true
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

export const ROUTES = {
  HOME: 'home',
  BLOCK: 'block',
  ADDRESS: 'addr',
  ADDRESS_FULL: 'address',
  RICHLIST: 'richlist',
  TRANSACTION: 'tx',
  SETTINGS: 'settings',
  VERIFY: 'verify',
};

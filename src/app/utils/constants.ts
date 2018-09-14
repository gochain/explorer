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
  {
    title: 'Wallet',
    link: 'https://wallet.gochain.io',
    icon: 'fa fa-wallet fa-fw',
    external: true
  },
  /*{
    title: 'Settings',
    link: '/settings',
    icon: 'fa fa-cogs fa-fw',
  },*/
];

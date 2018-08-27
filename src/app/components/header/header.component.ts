import {Component} from '@angular/core';
import {MenuItem} from '../../models/menu_item.model';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent {
  navItems: MenuItem[] = [
    {
      title: 'Blocks',
      link: '/home',
      icon: 'fa fa-link'
    },
    {
      title: 'Rich List',
      link: '/richlist',
      icon: 'fa fa-bars'
    },
    {
      title: 'Wallet',
      link: 'https://wallet.gochain.io',
      icon: 'fa fa-wallet',
      external: true
    },
  ];
}

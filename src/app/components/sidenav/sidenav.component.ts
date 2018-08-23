import {Component, OnInit} from '@angular/core';
import {MenuItem} from '../../models/menu_item.model';
import {LayoutService} from '../../services/template.service';

@Component({
  selector: 'app-sidenav',
  templateUrl: './sidenav.component.html',
  styleUrls: ['./sidenav.component.scss']
})
export class SidenavComponent implements OnInit {
  isOpen = true;
  navItems: MenuItem[] = [
    {
      title: 'Blocks',
      link: '/',
      icon: 'polymer'
    },
    {
      title: 'Rich List',
      link: '/richlist',
      icon: 'list'
    },
    {
      title: 'Wallet',
      link: 'https://wallet.gochain.io',
      icon: 'account_balance_wallet'
    },
  ];

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isSidenavOpen.subscribe((state: boolean) => {
      this.isOpen = state;
    });
  }

  toggle() {
    this._layoutService.toggleSidenav();
  }
}

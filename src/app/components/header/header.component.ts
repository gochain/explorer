import {Component, OnInit} from '@angular/core';
import {MenuItem} from '../../models/menu_item.model';
import {LayoutService} from '../../services/template.service';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss']
})
export class HeaderComponent implements OnInit {
  themeColor: string;

  navItems: MenuItem[] = [
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
    {
      title: 'Settings',
      link: '/settings',
      icon: 'fa fa-cogs fa-fw',
    },
  ];

  constructor(private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.themeColor.subscribe(value => {
      this.themeColor = value;
    })
  }
}

/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICES*/
import {MetaService} from '../../../services/meta.service';
/*UTILS*/
import {META_TITLES} from '../../../utils/constants';

@Component({
  selector: 'app-wallet-main',
  templateUrl: './wallet-main.component.html',
  styleUrls: ['./wallet-main.component.scss']
})
export class WalletMainComponent implements OnInit {

  constructor(
    private metaService: MetaService,
  ) {
  }

  ngOnInit() {
    this.metaService.setTitle(META_TITLES.WALLET.title);
  }
}

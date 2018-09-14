import { Component, OnInit } from '@angular/core';
import {LayoutService} from '../../services/layout.service';

@Component({
  selector: 'app-mobile-header',
  templateUrl: './mobile-header.component.html',
  styleUrls: ['./mobile-header.component.scss']
})
export class MobileHeaderComponent implements OnInit {

  constructor(private _layoutService: LayoutService) { }

  ngOnInit() {
  }

  toggleMenu() {
    this._layoutService.mobileMenuState.next(true);
  }
}

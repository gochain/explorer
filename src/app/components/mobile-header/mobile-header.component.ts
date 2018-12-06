/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICES*/
import {LayoutService} from '../../services/layout.service';

@Component({
  selector: 'app-mobile-header',
  templateUrl: './mobile-header.component.html',
  styleUrls: ['./mobile-header.component.scss']
})
export class MobileHeaderComponent implements OnInit {

  themeColor: string;

  constructor(public layoutService: LayoutService) {
  }

  ngOnInit() {
    this.layoutService.themeColor.subscribe(value => {
      this.themeColor = value;
    });
  }

  toggleMenu() {
    this.layoutService.mobileMenuState.next(!this.layoutService.mobileMenuState.value);
    this.layoutService.mobileSearchState.next(false);
  }

  toggleSearch() {
    this.layoutService.mobileSearchState.next(!this.layoutService.mobileSearchState.value);
    this.layoutService.mobileMenuState.next(false);
  }
}

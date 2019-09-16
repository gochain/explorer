import { Component, ElementRef } from '@angular/core';

@Component({
  selector: 'tab-title',
  templateUrl: 'tab-title.component.html'
})
export class TabTitleComponent {
  constructor(public element: ElementRef) { }
}

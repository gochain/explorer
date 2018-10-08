import { Component, ElementRef } from '@angular/core';

@Component({
  selector: 'tab-content',
  templateUrl: 'tab-content.component.html'
})
export class TabContentComponent {
  constructor(public element: ElementRef) { }
}

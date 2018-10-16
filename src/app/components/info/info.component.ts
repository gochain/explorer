import {Component, Input} from '@angular/core';

enum ContentTypes {
  ERROR = 'error',
  INFO = 'info',
  WARNING = 'warning',
  SUCCESS = 'success',
}

@Component({
  selector: 'app-info',
  templateUrl: './info.component.html',
  styleUrls: ['./info.component.scss']
})
export class InfoComponent {
  @Input() type?: ContentTypes;
  @Input() text: string;
}

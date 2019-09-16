import { Component } from '@angular/core';
import { ToastrService } from './toastr.service';

@Component({
  selector: 'app-toastr',
  templateUrl: './toastr.component.html',
  styleUrls: ['./toastr.component.scss']
})
export class ToastrComponent {
  constructor(public toastrService: ToastrService) { }
}

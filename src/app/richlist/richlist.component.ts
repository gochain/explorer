import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../api.service';
import { RichList } from "../rich_list";
import { of } from 'rxjs';
import { catchError, tap, map } from 'rxjs/operators';


@Component({
  selector: 'app-richlist',
  templateUrl: './richlist.component.html',
  styleUrls: ['./richlist.component.css']
})
export class RichlistComponent implements OnInit {

  richList: RichList;
  constructor(private api: ApiService) {
  }
  ngOnInit() {
    this.api.getRichlist(0, 100).subscribe((data: RichList) => {
      this.richList = data;
    }
    );

  }
}

/*CORE*/
import {Component, OnInit} from '@angular/core';
import {CommonService} from '../../services/common.service';
import {BlockList} from '../../models/block_list.model';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  recentBlocks: BlockList;

  constructor(private _commonService: CommonService) {
  }

  ngOnInit() {
    this._commonService.getRecentBlocks().subscribe((data: BlockList) => {
        this.recentBlocks = data;
      }
    );
  }
}

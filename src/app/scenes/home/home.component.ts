/*CORE*/
import {Component, OnInit} from '@angular/core';
/*SERVICES*/
import {LayoutService} from '../../services/template.service';
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {BlockList} from '../../models/block_list.model';


@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  recentBlocks: BlockList;

  constructor(private _commonService: CommonService, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
    this._commonService.getRecentBlocks().subscribe((data: BlockList) => {
        this.recentBlocks = data;
        this._layoutService.isPageLoading.next(false);
      }
    );
  }
}

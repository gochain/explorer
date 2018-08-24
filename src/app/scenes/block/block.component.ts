/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {ActivatedRoute, ParamMap} from '@angular/router';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/template.service';
/*MODELS*/
import {Block} from '../../models/block.model';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.css']
})
export class BlockComponent implements OnInit {

  private _blockNum: number;
  block: Observable<Block>;

  constructor(private _commonService: CommonService, private _route: ActivatedRoute, private _layoutService: LayoutService) {
  }

  ngOnInit() {
    this._layoutService.isPageLoading.next(true);
    this.block = this._route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this._blockNum = +params.get('id');
        return this._commonService.getBlock(this._blockNum);
      })
    );
    this._layoutService.isPageLoading.next(false);
  }

}

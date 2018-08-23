/*CORE*/
import {Component, OnInit} from '@angular/core';
import {Observable} from 'rxjs';
import {switchMap} from 'rxjs/operators';
import {ActivatedRoute, ParamMap} from '@angular/router';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
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

  constructor(private _commonService: CommonService, private route: ActivatedRoute) {
  }

  ngOnInit() {
    this.block = this.route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this._blockNum = +params.get('id');
        return this._commonService.getBlock(this._blockNum);
      })
    );
  }

}

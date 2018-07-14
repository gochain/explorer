import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { switchMap, tap } from 'rxjs/operators';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Block } from "../block";
import { ApiService } from '../api.service';

@Component({
  selector: 'app-block',
  templateUrl: './block.component.html',
  styleUrls: ['./block.component.css']
})
export class BlockComponent implements OnInit {

  private  blockNum: number;
  block: Observable<Block>;
  // block: Block;

  constructor(private api: ApiService, private route: ActivatedRoute) { }

  ngOnInit() {
    // this.route.paramMap.pipe(
    //   switchMap((params: ParamMap) => {
    //     // this is pretty jacked, probably a better way to do this, the examples don't subscribe.
    //     // https://angular.io/guide/router
    //     this.blockNum = +params.get('id');
    //     console.log("blocknum", this.blockNum);
    //     let ob = this.api.getBlock(this.blockNum);
    //     ob.subscribe((data: Block) => {
    //       console.log("block", data);
    //       this.block = data;
    //     })
    //     return ob;
    //   })
    // )
     this.block = this.route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this.blockNum = +params.get('id');
        console.log("blocknum", this.blockNum);
        return this.api.getBlock(this.blockNum);
      })
    )
  }

}

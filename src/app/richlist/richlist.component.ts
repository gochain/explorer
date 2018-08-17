import { Component, OnInit } from '@angular/core';
import { ApiService } from '../api.service';
import { RichList } from "../rich_list";


@Component({
  selector: 'app-richlist',
  templateUrl: './richlist.component.html',
  styleUrls: ['./richlist.component.css']
})
export class RichlistComponent implements OnInit {

  richList: RichList;
  
  skip = 0;
  limit = 100;
  isMoreDisabled = false;
  constructor(private api: ApiService) {
  }
  ngOnInit() {    
    this.richList = new RichList;
    this.richList.rankings = [];
    this.getMore();
  }
  getMore() {    
    this.api.getRichlist(this.skip, this.limit).subscribe((data: RichList) => {      
      this.richList.rankings = this.richList.rankings.concat(data.rankings);
      this.richList.circulating_supply = data.circulating_supply;
      this.richList.total_supply = data.total_supply;
      this.skip += 100;
      if (data.rankings.length < this.limit) {        
        this.isMoreDisabled = true;
      }
      this.richList.rankings = this.richList.rankings.map(acct => {
        acct.supply_owned = (acct.balance / this.richList.circulating_supply * 100).toFixed(2);
        return acct;
    });
    });
  }  

}
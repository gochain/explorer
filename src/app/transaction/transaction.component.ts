import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { switchMap, tap } from 'rxjs/operators';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Transaction } from "../transaction";
import { ApiService } from '../api.service';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.css']
})
export class TransactionComponent implements OnInit {

  private  txHash: string;
  transaction: Observable<Transaction>;

  constructor(private api: ApiService, private route: ActivatedRoute) { }

  ngOnInit() {    
     this.transaction = this.route.paramMap.pipe(
      switchMap((params: ParamMap) => {
        this.txHash = params.get('id');
        console.log("txnum", this.txHash);
        return this.api.getTransaction(this.txHash);
      })
    )
  }

}

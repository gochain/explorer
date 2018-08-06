import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { switchMap, tap } from 'rxjs/operators';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { Address } from "../address";
import { Transaction } from "../transaction";
import { ApiService } from '../api.service';

@Component({
  selector: 'app-address',
  templateUrl: './address.component.html',
  styleUrls: ['./address.component.css']
})
export class AddressComponent implements OnInit {

  private addrHash: string;
  address: Observable<Address>;
  transactions: Observable<Transaction[]>;

  constructor(private api: ApiService, private route: ActivatedRoute) { }

  ngOnInit() {    
    this.address = this.route.paramMap.pipe(
      switchMap((params: ParamMap) => {        
        this.addrHash = params.get('id');
        console.log("addrhash", this.addrHash);
        this.transactions = this.api.getAddressTransactions(this.addrHash);
        return this.api.getAddress(this.addrHash);
      })
    )
  }

}

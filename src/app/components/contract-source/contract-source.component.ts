/*CORE*/
import {Component, Input, OnInit} from '@angular/core';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
/*MODELS*/
import {Contract} from '../../models/contract.model';
import {Address} from '../../models/address.model';

@Component({
  selector: 'app-contract-source',
  templateUrl: './contract-source.component.html',
  styleUrls: ['./contract-source.component.css']
})
export class ContractSourceComponent implements OnInit {
  @Input()
  set addr(value: Address) {
    this._addr = value;
    this.getContractData();
  }

  get addr(): Address {
    return this._addr;
  }

  private _addr: Address;

  contract: Contract;

  constructor(
    private _commonService: CommonService,
  ) {
  }

  ngOnInit() {
  }

  getContractData() {
    this._commonService.getContract(this.addr.address).subscribe((data: Contract) => {
      this.contract = data;
    });
  }
}

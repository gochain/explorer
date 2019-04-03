/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, Params} from '@angular/router';
import {Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../../modules/wallet/wallet.service';
import {ApiService} from '../../services/api.service';
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {ABIDefinition} from 'web3/eth/abi';
import {TokenMetadata} from '../../models/token-metadata';

const TOKEN_URL_ABI: ABIDefinition = {
  'constant': true,
  'inputs': [{'name': '_tokenId', 'type': 'uint256'}],
  'name': 'tokenURI',
  'outputs': [{'name': '', 'type': 'string'}],
  'payable': false,
  'stateMutability': 'view',
  'type': 'function'
};

@Component({
  selector: 'app-token-asset',
  templateUrl: './token-asset.component.html',
  styleUrls: ['./token-asset.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class TokenAssetComponent implements OnInit, OnDestroy {
  contractAddr: string;
  tokenId: string;
  metadata: TokenMetadata;
  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService,
              private _route: ActivatedRoute,
              private _layoutService: LayoutService,
              private _walletService: WalletService,
              private _apiService: ApiService,
              private _toastrService: ToastrService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._route.params.pipe(
        filter((params: Params) => !!params.id && !!params.tokenId),
      ).subscribe((params: Params) => {
        this.contractAddr = params.id;
        this.tokenId = params.tokenId;
        this.metadata = null;
        this._layoutService.onLoading();
        this.getData();
      })
    );
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  getData() {
    let funcABI: string;
    try {
      funcABI = this._walletService.w3.eth.abi.encodeFunctionCall(TOKEN_URL_ABI, [this.tokenId]);
    } catch (err) {
      this._toastrService.danger(err);
      return;
    }
    this._walletService.w3.eth.call({
      to: this.contractAddr,
      data: funcABI,
    }).then((result: string) => {
      if (!result) {
        this._layoutService.offLoading();
        return;
      }
      const decoded: object = this._walletService.w3.eth.abi.decodeLog(TOKEN_URL_ABI.outputs, result, []);
      if (!decoded || !decoded[0]) {
        this._layoutService.offLoading();
        return;
      }
      const url: string = decoded[0];
      this._apiService.get(url, null, true).subscribe(res => {
        this._layoutService.offLoading();
        if (!res) {
          return;
        }
        const metadata = new TokenMetadata();
        metadata.name = res.name || null;
        metadata.description = res.description || null;
        metadata.image = res.image || null;
        metadata.external_url = res.external_url || null;
        metadata.origin_data = JSON.stringify(res, null, 4);
        this.metadata = metadata;
      });
    }).catch(err => {
      this._layoutService.offLoading();
      this._toastrService.danger(err);
    });

    // checking if contract is erc721metadata
    /*this._commonService.getAddress(this.contractAddr).pipe(
      filter((value: Address) => {
        if (!value || !value.contract || !value.erc_types.includes('Erc721Metadata')) {
          this._layoutService.offLoading();
          return false;
        }
        return true;
      }),
    ).subscribe((value: Address) => {
    });*/
  }
}

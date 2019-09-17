/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, Params} from '@angular/router';
import {Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
/*SERVICES*/
import {WalletService} from '../../services/wallet.service';
import {ApiService} from '../../services/api.service';
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {TokenMetadata} from '../../models/token-metadata';
import {AbiItem} from 'web3-utils';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES} from '../../utils/constants';

const TOKEN_URL_ABI: AbiItem = {
  'constant': true,
  'inputs': [{'name': '_tokenId', 'type': 'uint256'}],
  'name': 'tokenURI',
  'outputs': [{'name': '', 'type': 'string'}],
  'payable': false,
  'stateMutability': 'view',
  'type': 'function'
};

const OWNER_OF_ABI: AbiItem = {
  'constant': true,
  'inputs': [
    {
      'name': 'tokenId',
      'type': 'uint256'
    }
  ],
  'name': 'ownerOf',
  'outputs': [
    {
      'name': 'owner',
      'type': 'address'
    }
  ],
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
              private _metaService: MetaService,
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
        this.getData();
      })
    );
    this._metaService.setTitle(META_TITLES.TOKEN.title);
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  getData() {
    this._layoutService.onLoading();
    Promise.all([
      this._walletService.call(this.contractAddr, TOKEN_URL_ABI, [this.tokenId]),
      this._walletService.call(this.contractAddr, OWNER_OF_ABI, [this.tokenId]),
    ]).then(([tokenUrl, ownerData]: [object, object]) => {
      const url: string = tokenUrl[0];
      const metadata = new TokenMetadata();
      metadata.ownerAddr = ownerData['owner'];
      this._apiService.get(url, null, true).subscribe((res: any) => {
        metadata.name = res.name || null;
        if (metadata.name) {
          this._metaService.setTitle(`${META_TITLES.TOKEN.title} ${metadata.name}`);
        }
        metadata.description = res.description || null;
        metadata.image = res.image || null;
        metadata.external_url = res.external_url || null;
        metadata.origin_data = JSON.stringify(res, null, 4);
        this.metadata = metadata;
      });

    }).catch(err => {
      this._toastrService.danger(err);
    }).then(() => {
      this._layoutService.offLoading();
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

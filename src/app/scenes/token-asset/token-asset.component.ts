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
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {TokenMetadata} from '../../models/token-metadata';
import {ABIDefinition} from 'web3/eth/abi';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES} from '../../utils/constants';

const TOKEN_URL_ABI: ABIDefinition = {
  'constant': true,
  'inputs': [{'name': '_tokenId', 'type': 'uint256'}],
  'name': 'tokenURI',
  'outputs': [{'name': '', 'type': 'string'}],
  'payable': false,
  'stateMutability': 'view',
  'type': 'function'
};

const OWNER_OF_ABI: ABIDefinition = {
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
      const res = JSON.parse('{"attributes":[{"trait_type":"base","value":"goldfish"},{"trait_type":"eyes","value":"big"},{"trait_type":"mouth","value":"surprised"},{"trait_type":"level","value":5},{"trait_type":"stamina","value":1.4},{"trait_type":"personality","value":"sad"},{"display_type":"boost_number","trait_type":"aqua_power","value":30},{"display_type":"boost_percentage","trait_type":"stamina_increase","value":15},{"display_type":"number","trait_type":"generation","value":2}],"description":"Friendly OpenSea Creature that enjoys long swims in the ocean.","external_url":"https://openseacreatures.io/5","image":"https://storage.googleapis.com/opensea-prod.appspot.com/creature/5.png","name":"Captain McCoy"}');
      /*this._apiService.get(url, null, true).subscribe((res: any) => {*/
        metadata.name = res.name || null;
        if (metadata.name) {
          this._metaService.setTitle(`${META_TITLES.TOKEN.title} ${metadata.name}`);
        }
        metadata.description = res.description || null;
        metadata.image = res.image || null;
        metadata.external_url = res.external_url || null;
        metadata.origin_data = JSON.stringify(res, null, 4);
        this.metadata = metadata;
      /*});*/

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

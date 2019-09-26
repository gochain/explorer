/*CORE*/
import {Injectable} from '@angular/core';
import {BehaviorSubject, forkJoin, Observable} from 'rxjs';
import {filter, flatMap, map, take, tap} from 'rxjs/operators';
import {Resolve} from '@angular/router';
/*SERVICES*/
import {ApiService} from './api.service';
/*MODELS*/
import {BlockList} from '../models/block_list.model';
import {Block} from '../models/block.model';
import {Transaction} from '../models/transaction.model';
import {Address} from '../models/address.model';
import {RichList} from '../models/rich_list.model';
import {Holder} from '../models/holder.model';
import {InternalTransaction} from '../models/internal-transaction.model';
import {Stats} from '../models/stats.model';
import {Contract} from '../models/contract.model';
import {SignerData, SignerStat} from '../models/signer-stats';
import {SignerNode} from '../models/signer-node';
/*UTILS*/
import {ContractAbi} from '../utils/types';
import {objIsEmpty} from '../utils/functions';

@Injectable()
export class CommonService implements Resolve<string> {
  contractsCache = {};

  private _rpcProvider$: BehaviorSubject<string>;

  get rpcProvider$(): Observable<string> {
    if (!this._rpcProvider$) {
      this._rpcProvider$ = new BehaviorSubject(null);
      return this.getRpcProvider().pipe(
        tap(value => this._rpcProvider$.next(value)),
      );
    }
    return this._rpcProvider$.pipe(
      filter(value => !!value),
      take(1),
    );
  }

  private _signers$: BehaviorSubject<any>;

  get signers$(): Observable<SignerNode> {
    if (!this._signers$) {
      this._signers$ = new BehaviorSubject<any>(null);
      this._apiService.get('/signers/list').subscribe(value => {
        this._signers$.next(value);
      });
    }
    return this._signers$.pipe(
      filter(value => !!value),
      take(1),
    );
  }

  constructor(private _apiService: ApiService) {
  }

  /**
   * getting RpcProvider
   */
  resolve(): Observable<string> {
    return this.rpcProvider$;
  }

  getRpcProvider(): Observable<string> {
    return this._apiService.get('/rpc_provider');
  }

  getAbi(): Observable<ContractAbi> {
    return this._apiService.get('/assets/data/abi.json', null, true);
  }

  getApiUrl(): string {
    return this._apiService.apiURL;
  }

  getRecentBlocks(): Observable<BlockList> {
    return this._apiService.get('/blocks');
  }

  getBlock(blockNum: number | string, data?: any): Observable<Block> {
    return forkJoin([
      this.signers$,
      this._apiService.get('/blocks/' + blockNum, data),
    ]).pipe(
      map(([signers, block]: [SignerNode, Block]) => {
        block.signerDetails = signers[block.miner.toLowerCase()] || null;
        if (block.extra && block.extra.candidate) {
          block.extra.signerDetails = signers[block.extra.candidate.toLowerCase()] || null;
        }
        return block;
      }),
    );
  }

  checkBlockExist(blockHash: string) {
    return this._apiService.head('/blocks/' + blockHash);
  }

  checkTransactionExist(blockHash: string) {
    return this._apiService.head('/transaction/' + blockHash);
  }


  getBlockTransactions(blockNum: number | string, data?: any) {
    return this._apiService.get('/blocks/' + blockNum + '/transactions', data);
  }

  getTransaction(txHash: string): Observable<Transaction | null> {
    return this._apiService.get('/transaction/' + txHash);
  }

  getAddress(addrHash: string): Observable<Address> {
    return forkJoin([
      this.signers$,
      this._apiService.get('/address/' + addrHash),
    ]).pipe(
      map(([signers, address]: [SignerNode, Address]) => {
        address.signerDetails = signers[address.address.toLowerCase()] || null;
        return address;
      })
    );
  }

  getAddressTransactions(addrHash: string, data?: any): Observable<Transaction[]> {
    return this._apiService.get('/address/' + addrHash + '/transactions', data);
  }

  getAddressHolders(addrHash: string, data?: any): Observable<Holder[]> {
    return this._apiService.get('/address/' + addrHash + '/holders', data);
  }

  getAddressTokens(addrHash: string, data?: any): Observable<any> {
    return this._apiService.get(`/address/${addrHash}/owned_tokens`, data);
  }

  getAddressInternalTransaction(addrHash: string, data?: any): Observable<InternalTransaction[]> {
    return this._apiService.get('/address/' + addrHash + '/internal_transactions', data);
  }

  getContract(addrHash: string): Observable<Contract> {
    return this._apiService.get('/address/' + addrHash + '/contract');
  }

  getRichlist(data?: any): Observable<RichList> {
    return this._apiService.get('/richlist', data);
  }

  getContractsList(data?: any): Observable<Address[]> {
    return this._apiService.get('/contracts', data).pipe(
      map(v => v || []),
    );
  }

  getStats(): Observable<Stats> {
    return this._apiService.get('/stats');
  }

  getSignerStats(): Observable<SignerStat[]> {
    return this.signers$.pipe(
      flatMap((signers: SignerNode) => {
        return this._apiService.get('/signers/stats').pipe(
          tap((stats: SignerStat[]) => {
            if (objIsEmpty(signers)) {
              return;
            }
            stats.forEach((stat: SignerStat) => {
              stat.signer_stats.forEach((signer: SignerData) => {
                signer.data = signers[signer.signer_address];
              });
            });
          }),
        );
      }),
    );
  }
}

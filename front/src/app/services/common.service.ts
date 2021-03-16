/*CORE*/
import { Injectable } from '@angular/core';
import { BehaviorSubject, forkJoin, Observable } from 'rxjs';
import { filter, flatMap, map, take, tap } from 'rxjs/operators';
import { Resolve } from '@angular/router';
/*SERVICES*/
import { ApiService } from './api.service';
/*MODELS*/
import { BlockList } from '../models/block_list.model';
import { Block } from '../models/block.model';
import { Transaction } from '../models/transaction.model';
import { Address } from '../models/address.model';
import { RichList } from '../models/rich_list.model';
import { Holder } from '../models/holder.model';
import { InternalTransaction } from '../models/internal-transaction.model';
import { Stats } from '../models/stats.model';
import { SupplyStats } from '../models/supply.model';
import { Contract } from '../models/contract.model';
import { SignerData, SignerStat } from '../models/signer-stats';
import { SignerNode } from '../models/signer-node';
/*UTILS*/
import { ContractAbi, ContractEventsAbi, ContractAbiByID, AbiItemIDed } from '../utils/types';
import { FunctionName } from '../utils/enums';
import { objIsEmpty } from '../utils/functions';
import { AbiItem } from 'web3-utils';

@Injectable()
export class CommonService implements Resolve<string> {
  contractsCache = {};

  private _rpcProvider$: BehaviorSubject<string>;
  get rpcProvider$(): Observable<string> {
    if (!this._rpcProvider$) {
      this._rpcProvider$ = new BehaviorSubject(null);
      this.getRpcProvider().subscribe(v => {
        this._rpcProvider$.next(v);
      });
    }
    return this._rpcProvider$.pipe(
      filter(v => !!v),
      take(1),
    );
  }

  private _abi$: BehaviorSubject<ContractAbi>;
  get abi$() {
    if (!this._abi$) {
      this.initAbi();
    }
    return this._abi$.pipe(
      filter(v => v !== null),
      take(1),
    );
  }

  private _abiByID$: BehaviorSubject<ContractAbiByID>;
  get abiByID$(): Observable<ContractAbiByID> {
    if (!this._abiByID$) {
      this.initAbi();
    }
    return this._abiByID$.pipe(
      filter(v => v !== null),
      take(1),
    );
  }

  private _eventsAbi$: BehaviorSubject<ContractEventsAbi>;
  get eventsAbi$(): Observable<ContractEventsAbi> {
    if (!this._eventsAbi$) {
      this._eventsAbi$ = new BehaviorSubject<ContractEventsAbi>(null);
      this.getEventsAbi().subscribe(v => {
        this._eventsAbi$.next(v);
      });
    }
    return this._eventsAbi$.pipe(
      filter(v => v !== null),
      take(1),
    );
  }

  private _signers$: BehaviorSubject<any>;
  get signers$(): Observable<SignerNode> {
    if (!this._signers$) {
      this._signers$ = new BehaviorSubject<any>(null);
      this._apiService.get('/signers/list').subscribe(v => {
        this._signers$.next(v);
      });
    }
    return this._signers$.pipe(
      filter(v => !!v),
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

  getFunctionsAbi(): Observable<ContractAbi> {
    return this._apiService.get('/assets/abi/functions.json', null, true);
  }

  getEventsAbi(): Observable<ContractEventsAbi> {
    return this._apiService.get('/assets/abi/events.json', null, true);
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

  getTransaction(hash: string, nonceId: string): Observable<Transaction | null> {
    if (nonceId) {
      return this._apiService.get('/address/' + hash + '/tx/' + nonceId);
    } else {
      return this._apiService.get('/transaction/' + hash);
    }
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

  getAddressTokens(addrHash: string, data?: any): Observable<Holder> {
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

  getSupplyStats(): Observable<SupplyStats> {
    return this._apiService.get('/supply');
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

  private initAbi() {
    this._abi$ = new BehaviorSubject<ContractAbi>(null);
    this._abiByID$ = new BehaviorSubject<ContractAbiByID>(null);
    this.getFunctionsAbi().subscribe(v => {
      const abi: ContractAbi = <ContractAbi>{};
      const abiByID: ContractAbiByID = {};
      Object.entries(v).forEach((value: [FunctionName, AbiItemIDed]) => {
        abi[value[0]] = <AbiItem>value[1];
        abiByID[value[1].id] = <AbiItem>value[1];
      });
      this._abi$.next(abi);
      this._abiByID$.next(abiByID);
    });
  }
}

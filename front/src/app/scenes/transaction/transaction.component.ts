/*CORE*/
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap} from '@angular/router';
import {forkJoin, interval, Observable, of, Subscription} from 'rxjs';
import {catchError, map, mergeMap, startWith, tap} from 'rxjs/operators';
import {fromPromise} from 'rxjs/internal-compatibility';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {LayoutService} from '../../services/layout.service';
import {WalletService} from '../../services/wallet.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {ProcessedLog, ProcessedABIItem, Transaction, TxLog, ProcessedABIData} from '../../models/transaction.model';
import {ContractAbiByID, ContractEventsAbi} from '../../utils/types';
import Web3 from 'web3';
import {AbiItem, AbiInput} from 'web3-utils';
import {Address} from '../../models/address.model';
import {Contract} from '../../models/contract.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES} from '../../utils/constants';
import {ErcName} from '../../utils/enums';
import CID from '../../../../node_modules/cids/dist/index.min.js';
import {convertWithDecimals} from '../../utils/functions';

@Component({
  selector: 'app-transaction',
  templateUrl: './transaction.component.html',
  styleUrls: ['./transaction.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class TransactionComponent implements OnInit, OnDestroy {

  showInputRaw = true; // Input data as raw plain text, or rich formatting.
  showUtf8 = false; // Raw input data as UTF8, or hex encoded bytes.
  showLogsRaw = false; // Logs as raw JSON, or rich formatting.
  tx: Transaction;

  recentBlockNumber$: Observable<number> = interval(5000).pipe(
    startWith(0),
    mergeMap(() => {
      return this._walletService.w3Call.pipe(mergeMap((web3: Web3) => {
        return fromPromise<number>(web3.eth.getBlockNumber());
      }));
    }),
  );

  private _subsArr$: Subscription[] = [];

  constructor(private _commonService: CommonService,
              private _route: ActivatedRoute,
              private _layoutService: LayoutService,
              private _walletService: WalletService,
              private _metaService: MetaService,
  ) {
  }

  ngOnInit() {
    this._layoutService.onLoading();
    this._subsArr$.push(
      this._route.paramMap.pipe(
        tap<ParamMap>(() => {
          this._layoutService.onLoading();
        }),
        mergeMap((params: ParamMap) => this.getTx(params)),
      ).subscribe((tx: (Transaction | null)) => {
        tx.input_data = '0x' + tx.input_data;
        tx.parsedLogs = JSON.parse(tx.logs);
        tx.prettifiedLogs = JSON.stringify(tx.parsedLogs, null, '\t');
        this.tx = tx;
        this.processTransaction(this.tx);
        this._layoutService.offLoading();
      })
    );
    this._metaService.setTitle(META_TITLES.TRANSACTION.title);
  }

  ngOnDestroy(): void {
    this._layoutService.offLoading();
  }

  private processTransaction(tx: Transaction): void {
    if (tx.input_data && tx.input_data !== '0x' && tx.input_data !== '0X') {
      let data: Observable<ProcessedABIData>;
      if (tx.to) {
        data = forkJoin([
          <ContractAbiByID>this._commonService.abiByID$,
          tx.to ? this._commonService.getAddress(tx.to) : of(null),
          tx.to ? this._commonService.getContract(tx.to) : of(null),
          this._walletService.w3Call,
        ]).pipe(
          map((res) => this.processTransactionInputTo(res)),
        );
      } else if (tx.contract_address) {
        data = this._commonService.getContract(tx.contract_address).pipe(
          map<Contract, ProcessedABIData>((contract: Contract) => this.processTransactionInputDeploy(contract))
        );
      }
      if (data) {
        data.subscribe((input: ProcessedABIData) => {
          this.showInputRaw = input == null;
          tx.processedInputData = input;
        });
      }
    }

    if (!!tx.parsedLogs.length) {
      forkJoin([
        this._commonService.eventsAbi$,
        this._walletService.w3Call,
      ]).pipe(
        mergeMap((res) => this.processTransactionLogs(res))
      ).subscribe(logs => {
        tx.processedLogs = logs;
      });
    }
  }

  private processTransactionInputTo([abis, address, contract, web3]: [ContractAbiByID, Address, Contract, Web3]): ProcessedABIData {
    const processedInputData = new ProcessedABIData();
    // Contract call.
    const methodId = this.tx.input_data.slice(0, 10);
    let c: AbiItem;
    if (contract && contract.abi) {
      // Check the attached verified abi.
      processedInputData.title = `${contract.contract_name}.`;
      // TODO this is inefficient - could be cached or precomputed
      c = contract.abi.find(a => methodId === web3.eth.abi.encodeFunctionSignature(a));
    }
    if (!c) {
      // Check known functions.
      processedInputData.title = '';
      c = abis[methodId];
    }
    if (!c) {
      // We don't recognize this method.
      return null;
    }
    processedInputData.title += c.name;
    try {
      const d: object = web3.eth.abi.decodeParameters(c.inputs, '0x' + this.tx.input_data.slice(10));
      processedInputData.items = c.inputs.map(this.processAbiItem(d, address));
    } catch (e) {
      console.error('failed to decode input data', e);
      processedInputData.items = c.inputs.map(input => {
        return <ProcessedABIItem>{name: input.name, value: '? error'};
      });
    }
    return processedInputData;
  }

  private processTransactionInputDeploy(contract: Contract): ProcessedABIData {
    const processedInputData = new ProcessedABIData();
    // Contract deploy.
    processedInputData.title = `new ${contract.contract_name}`;

    const c: AbiItem = contract.abi.find(a => a.type === 'constructor');
    if (c && c.inputs.length > 0) {
      // TODO we know how to decode the constructor args, and we know they are at the end of the input_data, but
      // we don't know where they begin.
      processedInputData.items = c.inputs.map((input: AbiInput) => {
        return <ProcessedABIItem>{name: input.name, value: '<unknown>'};
      });
    }
    return processedInputData;
  }

  private processTransactionLogs([events, web3]: [ContractEventsAbi, Web3]): Observable<ProcessedLog[]> {
    // Create one observable per unique address.
    const addrs: { string: Observable<Address> } = <{ string: Observable<Address> }>this.tx.parsedLogs.reduce((acc: object, log: TxLog) => {
      if (!acc[log.address] && !!log.topics.length && !!log.topics[0]) {
        acc[log.address] = this._commonService.getAddress(log.address);
      }
      return acc;
    }, {});
    // Fetch all the Addresses, then use them to parse the logs.
    return forkJoin(addrs).pipe(
      catchError(err => {
          console.error(`failed to get address: ${err}`);
          return of(null);
        }
      ),
      map((eventAddrs: Address[]) => {
        return this.tx.parsedLogs.map((log: TxLog) => {
          const processedLog: ProcessedLog = new ProcessedLog();
          processedLog.index = +log.logIndex;
          processedLog.contract_address = log.address;
          processedLog.removed = log.removed;
          processedLog.data = [];

          let abiItem: AbiItem;
          let decodedLog: object;
          const address: Address = eventAddrs[log.address];
          if (address) {
            if (!!log.topics.length && !!log.topics[0]) {
              const eventSignature = <string>log.topics[0];
              const knownEvent: object = events[eventSignature];
              if (knownEvent) {
                if (Object.keys(knownEvent).length === 1) {
                  // Only once choice.
                  abiItem = Object.values(knownEvent)[0];
                } else {
                  // Find a match.
                  const ercType: string = address.erc_types.find((item: string) => !!knownEvent[item]);
                  if (ercType) {
                    abiItem = knownEvent[ercType];
                  }
                }
                if (abiItem) {
                  try {
                    const data: string = !log.data || log.data === '0x' || log.data === '0X' ? null : log.data;
                    const topics: string[] = <string[]>log.topics.slice(1);
                    decodedLog = web3.eth.abi.decodeLog(abiItem.inputs, data, topics);
                  } catch (e) {
                    console.error('error occurred while decoding log', e);
                  }
                }
              }
            }
          }
          if (decodedLog) {
            const items: ProcessedABIItem[] = abiItem.inputs.map(this.processAbiItem(decodedLog, address));
            processedLog.data.push({
              title: abiItem.name,
              items,
            });
          } else {
            const items: ProcessedABIItem[] = log.topics.map(topic => (<ProcessedABIItem>{value: <string>topic}));
            if (log.topics.length) {
              processedLog.data.push(<ProcessedABIData>{
                title: 'topics',
                items,
              });
            }
            if (log.data && log.data !== '0x' && log.data !== '0X') {
              processedLog.data.push(<ProcessedABIData>{
                title: 'data',
                items: [{value: log.data}],
              });
            }
          }
          return processedLog;
        });
      }));
  }

  private processAbiItem(decoded: object, address: Address): (input: AbiInput) => ProcessedABIItem {
    return (input: AbiInput): ProcessedABIItem => {
      const item = new ProcessedABIItem();
      item.name = input.name;
      item.value = decoded[input.name];
      if (address.erc_types.includes(ErcName.Go20)) {
        if (input.name === 'value') {
          if (address.decimals) {
            const val = convertWithDecimals(decoded[input.name], false, true, address.decimals);
            item.value = `${val} ${address.token_symbol}`;
          }
        }
      }
      if (address.erc_types.includes(ErcName.Go721)) {
        if (input.name === 'tokenId') {
          item.link = `/token/${address.address}/asset/${decoded[input.name]}`;
          return item;
        } else if (input.name === 'tokenURI') {
          item.link = decoded[input.name];
          item.linkExternal = true;
          return item;
        }
      }
      if (input.name === 'cid' && input.type === 'bytes' && !input.indexed) {
        const hex = decoded[input.name];
        if (hex.length > 2) {
          try {
            const cid: CID = new CID(Buffer.from(hex.slice(2), 'hex'));
            const str = cid.toString();
            item.value = str;
            item.link = `https://ipfs.io/ipfs/${str}`;
            item.linkExternal = true;
            return item;
          } catch (e) {
            console.error(`failed to parse cid ${hex}: ${e}`);
          }
        }
      }
      // if (ercTypes.includes(ErcName.GoST)) {
      //   //TODO link ETH stuff
      // }
      if (input.type === 'address') {
        if (decoded[input.name] === '0x0000000000000000000000000000000000000000') {
          // Reformat empty addresses.
          switch (input.name) {
            case 'to':
              item.value = '0x0 (burn)';
              return item;
            case 'from':
              item.value = '0x0 (mint)';
              return item;
            case 'previousOwner':
            case 'newOwner':
              item.value = '0x0 (none)';
              return item;
            default:
              item.value = '0x0';
              return item;
          }
        }
        // Link non-empty addresses.
        item.link = `/addr/${decoded[input.name]}`;
        return item;
      }
      return item;
    };
  }

  /**
   * getting tx from server
   * @param hash
   * @param nonceId
   */
  private getTx(params: ParamMap): Observable<Transaction | null> {
    const hash = params.get('id');
    const nonceId = params.get('nonce_id');
    return this._commonService.getTransaction(hash, nonceId).pipe(
      mergeMap((tx: Transaction | null) => {
        if (!tx) {
          return this._walletService.getTxData(hash);
        }
        return of(tx);
      }),
    );
  }
}

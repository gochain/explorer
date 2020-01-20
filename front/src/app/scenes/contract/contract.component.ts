/*CORE*/
import {Component, OnInit} from '@angular/core';
import {ActivatedRoute, ParamMap, Router} from '@angular/router';
import {Observable, Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
/*SERVICES*/
import {ContractService} from '../../services/contract.service';
import {ToastrService} from '../../modules/toastr/toastr.service';
import {MetaService} from '../../services/meta.service';
/*MODELS*/
import {Contract} from '../../models/contract.model';
/*UTILS*/
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {META_TITLES, ROUTES} from '../../utils/constants';

// import {Compiler} from '../../models/compiler.model';

@Component({
  selector: 'app-contract',
  templateUrl: './contract.component.html',
  styleUrls: ['./contract.component.scss']
})
@AutoUnsubscribe('_subsArr$')
export class ContractComponent implements OnInit {
  contract: Contract;
  /*recaptchaPublicKey = environment.RECAPTCHA_KEY;*/

  form: FormGroup = this._fb.group({
    address: ['', [Validators.required, Validators.minLength(42), Validators.maxLength(42)]],
    contract_name: ['', Validators.required],
    compiler_version: ['', Validators.required],
    evm_version: [''],
    optimization: [true, Validators.required],
    source_code: ['', Validators.required],
    /*recaptcha_token: null,*/
  });

  compilers$: Observable<any[]> = this._contactService.getCompilersList();

  private _subsArr$: Subscription[] = [];

  constructor(private _activatedRoute: ActivatedRoute,
              private _fb: FormBuilder,
              private _contactService: ContractService,
              private _toastrService: ToastrService,
              private _router: Router,
              private _metaService: MetaService,
  ) {
  }

  ngOnInit() {
    this._subsArr$.push(
      this._activatedRoute.queryParamMap.pipe(
        filter((params: ParamMap) => params.has('address'))
      ).subscribe((params: ParamMap) => {
        const addr = params.get('address');
        if (addr.length === 42) {
          this.getContract(addr);
          this.form.patchValue({
            address: addr
          });
        } else {
          this._toastrService.warning('Contract address is invalid');
        }
      })
    );
    this._metaService.setTitle(META_TITLES.VERIFY.title);
  }

  getContract(addrHash: string) {
    this._subsArr$.push(this._contactService.getContract(addrHash).subscribe((contract: Contract) => {
      if (!contract) {
        this._toastrService.danger('Contract address not found');
      } else {
        this.contract = contract;
      }
    }));
  }

  onSubmit() {
    if (!this.form.valid) {
      this._toastrService.danger('Some field is not correct');
      return;
    }
    const data = this.form.getRawValue();
    this._contactService.compile(data).pipe(
      filter((contract: Contract) => !!contract)
    ).subscribe((contract: Contract) => {
      this.contract = contract;
      if (this.contract.valid) {
        this._toastrService.success('Contract has been successfully verified');
        this.form.reset();
        this._router.navigate(
          [`/${ROUTES.ADDRESS}/`, this.contract.address],
          {
            queryParams: {
              addr_tab: 'contract_source',
            },
          },
        );
      }
    });
  }
}

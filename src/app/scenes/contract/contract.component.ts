import {Component, OnInit} from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import {Subscription} from 'rxjs';
import {filter} from 'rxjs/operators';
import {Params} from '@angular/router/src/shared';
import {AutoUnsubscribe} from '../../decorators/auto-unsubscribe';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import {ContractService} from '../../services/contract.service';
import {Compiler} from '../../models/compiler.model';

@Component({
  selector: 'app-contract',
  templateUrl: './contract.component.html',
  styleUrls: ['./contract.component.css']
})
@AutoUnsubscribe('_subsArr$')
export class ContractComponent implements OnInit {
  compilers: any[] = [];

  form: FormGroup = this.fb.group({
    address: ['', Validators.required, Validators.minLength(42), Validators.maxLength(42)],
    contractName: [''],
    compilerVersion: ['', Validators.required],
    optimization: [false],
    sourceCode: ['', Validators.required],
    abi: [''],
  });

  private _subsArr$: Subscription[] = [];


  constructor(private _activatedRoute: ActivatedRoute, private fb: FormBuilder, private contactService: ContractService) {
    this.contactService.getCompilersList().subscribe((value: any) => {
      this.compilers = value.builds.map((item: Compiler) => {
        if (item.prerelease && item.prerelease.length > 0) {
          return item.version + '-' + item.prerelease;
        } else {
          return item.version;
        }
      });
    });
  }

  ngOnInit() {
    this._subsArr$.push(
      this._activatedRoute.queryParamMap.subscribe((params: ParamMap) => {
        const addr = params.get('address');
        if (addr.length === 42) {
          this.form.patchValue({
            address: addr
          });
        }
      })
    );
  }

  onSubmit() {
    const data = this.form.getRawValue();
    this.contactService.compile(data).subscribe(data => {
      console.log(data);
    });
  }
}

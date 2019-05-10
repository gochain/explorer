/*CORE*/
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
/*MODULES*/
import {AppModule} from '../../app.module';
/*COMPONENTS*/
import {TransactionComponent} from './transaction.component';
import {WalletCommonModule} from '../../modules/wallet/wallet-common.module';

describe('TransactionComponent', () => {
  let component: TransactionComponent;
  let fixture: ComponentFixture<TransactionComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, WalletCommonModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TransactionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

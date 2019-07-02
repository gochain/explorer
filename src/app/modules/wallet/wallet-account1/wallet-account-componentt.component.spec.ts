/*CORE*/
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {RouterTestingModule} from '@angular/router/testing';
/*MODULES*/
import {WalletModule} from '../wallet.module';
import {AppModule} from '../../../app.module';
/*COMPONENTS*/
import {WalletAccountComponentt} from './wallet-account-componentt.component';


describe('WalletUseComponent', () => {
  let component: WalletAccountComponentt;
  let fixture: ComponentFixture<WalletAccountComponentt>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [RouterTestingModule, AppModule, WalletModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletAccountComponentt);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

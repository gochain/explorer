/*CORE*/
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {RouterTestingModule} from '@angular/router/testing';
/*MODULES*/
import {WalletModule} from '../wallet.module';
import {AppModule} from '../../../app.module';
/*COMPONENTS*/
import {WalletUseComponent} from './wallet-use.component';


describe('WalletUseComponent', () => {
  let component: WalletUseComponent;
  let fixture: ComponentFixture<WalletUseComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [RouterTestingModule, AppModule, WalletModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletUseComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

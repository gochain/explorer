import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {WalletUseComponent} from './wallet-use.component';
import {WalletModule} from '../wallet.module';
import {RouterTestingModule} from '@angular/router/testing';

describe('WalletUseComponent', () => {
  let component: WalletUseComponent;
  let fixture: ComponentFixture<WalletUseComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [RouterTestingModule, WalletModule]
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

import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {WalletSendComponent} from './wallet-send.component';
import {WalletModule} from '../wallet.module';
import {AppModule} from '../../../app.module';

describe('WalletSendComponent', () => {
  let component: WalletSendComponent;
  let fixture: ComponentFixture<WalletSendComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule, WalletModule],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletSendComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

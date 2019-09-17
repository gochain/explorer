import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {WalletAccountComponent} from './wallet-account.component';
import {AppModule} from '../../app.module';

describe('WalletAccountComponent', () => {
  let component: WalletAccountComponent;
  let fixture: ComponentFixture<WalletAccountComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletAccountComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

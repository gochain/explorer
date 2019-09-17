import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {WalletMainComponent} from './wallet-main.component';
import {AppModule} from '../../app.module';

describe('WalletMainComponent', () => {
  let component: WalletMainComponent;
  let fixture: ComponentFixture<WalletMainComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletMainComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

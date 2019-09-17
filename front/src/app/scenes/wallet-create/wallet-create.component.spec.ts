import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {WalletCreateComponent} from './wallet-create.component';
import {AppModule} from '../../app.module';

describe('WalletCreateComponent', () => {
  let component: WalletCreateComponent;
  let fixture: ComponentFixture<WalletCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule],
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(WalletCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

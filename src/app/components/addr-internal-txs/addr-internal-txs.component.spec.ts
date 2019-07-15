import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddrInternalTxsComponent } from './addr-internal-txs.component';
import {AppModule} from '../../app.module';

describe('AddrInternalTxsComponent', () => {
  let component: AddrInternalTxsComponent;
  let fixture: ComponentFixture<AddrInternalTxsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddrInternalTxsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddrInternalTxsComponent } from './addr-internal-txs.component';

describe('AddrInternalTxsComponent', () => {
  let component: AddrInternalTxsComponent;
  let fixture: ComponentFixture<AddrInternalTxsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AddrInternalTxsComponent ]
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

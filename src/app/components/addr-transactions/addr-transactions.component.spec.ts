import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddrTransactionsComponent } from './addr-transactions.component';
import {AppModule} from '../../app.module';

describe('AddrTransactionsComponent', () => {
  let component: AddrTransactionsComponent;
  let fixture: ComponentFixture<AddrTransactionsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddrTransactionsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

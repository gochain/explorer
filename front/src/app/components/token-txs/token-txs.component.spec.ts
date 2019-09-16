import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TokenTxsComponent } from './token-txs.component';
import {AppModule} from '../../app.module';

describe('TokenTxsComponent', () => {
  let component: TokenTxsComponent;
  let fixture: ComponentFixture<TokenTxsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TokenTxsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

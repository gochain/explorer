import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {TokenHoldersComponent} from './token-holders.component';
import {AppModule} from '../../app.module';

describe('TokenHoldersComponent', () => {
  let component: TokenHoldersComponent;
  let fixture: ComponentFixture<TokenHoldersComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TokenHoldersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

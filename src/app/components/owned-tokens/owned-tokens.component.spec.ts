import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {OwnedTokensComponent} from './owned-tokens.component';
import {AppModule} from '../../app.module';

describe('OwnedTokensComponent', () => {
  let component: OwnedTokensComponent;
  let fixture: ComponentFixture<OwnedTokensComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OwnedTokensComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

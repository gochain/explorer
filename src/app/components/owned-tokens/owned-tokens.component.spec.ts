import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OwnedTokensComponent } from './owned-tokens.component';

describe('OwnedTokensComponent', () => {
  let component: OwnedTokensComponent;
  let fixture: ComponentFixture<OwnedTokensComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OwnedTokensComponent ]
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

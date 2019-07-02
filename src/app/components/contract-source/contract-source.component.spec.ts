import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ContractSourceComponent } from './contract-source.component';

describe('ContractSourceComponent', () => {
  let component: ContractSourceComponent;
  let fixture: ComponentFixture<ContractSourceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ContractSourceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContractSourceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

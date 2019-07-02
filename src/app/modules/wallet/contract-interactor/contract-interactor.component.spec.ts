import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ContractInteractorComponent } from './contract-interactor.component';

describe('ContractInteractorComponent', () => {
  let component: ContractInteractorComponent;
  let fixture: ComponentFixture<ContractInteractorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ContractInteractorComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ContractInteractorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

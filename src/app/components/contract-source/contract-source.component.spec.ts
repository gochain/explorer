import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ContractSourceComponent} from './contract-source.component';
import {AppModule} from '../../app.module';

describe('ContractSourceComponent', () => {
  let component: ContractSourceComponent;
  let fixture: ComponentFixture<ContractSourceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
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

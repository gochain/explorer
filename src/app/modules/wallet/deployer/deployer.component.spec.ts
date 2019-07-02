import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DeployerComponent } from './deployer.component';

describe('DeployerComponent', () => {
  let component: DeployerComponent;
  let fixture: ComponentFixture<DeployerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DeployerComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DeployerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

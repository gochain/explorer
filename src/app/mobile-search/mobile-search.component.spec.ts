import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { MobileSearchComponent } from './mobile-search.component';

describe('MobileSearchComponent', () => {
  let component: MobileSearchComponent;
  let fixture: ComponentFixture<MobileSearchComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ MobileSearchComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MobileSearchComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

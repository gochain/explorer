import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RichlistComponent } from './richlist.component';

describe('RichlistComponent', () => {
  let component: RichlistComponent;
  let fixture: ComponentFixture<RichlistComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RichlistComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RichlistComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

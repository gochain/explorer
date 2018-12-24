import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SendTxComponent } from './send-tx.component';

describe('SendTxComponent', () => {
  let component: SendTxComponent;
  let fixture: ComponentFixture<SendTxComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SendTxComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SendTxComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

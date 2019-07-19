import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {SignersComponent} from './signers.component';
import {AppModule} from '../../app.module';

describe('SignersComponent', () => {
  let component: SignersComponent;
  let fixture: ComponentFixture<SignersComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SignersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

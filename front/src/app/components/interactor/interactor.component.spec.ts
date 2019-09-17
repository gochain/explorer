import {async, ComponentFixture, TestBed} from '@angular/core/testing';
import {InteractorComponent} from './interactor.component';
import {AppModule} from '../../app.module';

describe('ContractInteractorComponent', () => {
  let component: InteractorComponent;
  let fixture: ComponentFixture<InteractorComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(InteractorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

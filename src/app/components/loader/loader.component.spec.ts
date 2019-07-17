/*
/!*CORE*!/
import {async, ComponentFixture, TestBed} from '@angular/core/testing';
/!*MODULES*!/
import {AppModule} from '../../app.module';
/!*COMPONENTS*!/
import {LoaderComponent} from './loader.component';

describe('LoaderComponent', () => {
  let component: LoaderComponent;
  let fixture: ComponentFixture<LoaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LoaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
*/

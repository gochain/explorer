/*CORE*/
import {async, TestBed} from '@angular/core/testing';
/*MODULES*/
import {AppModule} from './app.module';
/*COMPONENTS*/
import {AppComponent} from './app.component';

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule,
      ],
    }).compileComponents();
  }));
  it('should create the app', async(() => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  }));
  it(`should have as isPageLoading 'false'`, async(() => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.isPageLoading).toEqual(false);
  }));
});

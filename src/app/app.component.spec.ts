/*CORE*/
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
/*MODULES*/
import { AppModule } from './app.module';
/*COMPONENTS*/
import { AppComponent } from './app.component';

let comp: AppComponent;
let fixture: ComponentFixture<AppComponent>;

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule,
      ],
    }).compileComponents();
  }));
  beforeEach(() => {
    fixture = TestBed.createComponent(AppComponent);
    fixture.detectChanges();
    comp = fixture.debugElement.componentInstance;
  });
  /**
   * BEFORE INIT
   */
  it(`should be initialized`, () => {
    expect(fixture).toBeDefined();
    expect(comp).toBeDefined();
  });
  /**
   * INIT
   */
  it('should create the app', async(() => {
    expect(comp).toBeTruthy();
  }));
  /**
   * DEFAULT VALUES
   */
  it(`should have as isPageLoading 'false'`, async(() => {
    expect(comp.isPageLoading).toEqual(false);
  }));
});

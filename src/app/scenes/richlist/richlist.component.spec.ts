/*CORE*/
import { async, ComponentFixture, TestBed } from '@angular/core/testing';
/*MODULES*/
import {AppModule} from '../../app.module';
/*COMPONENTS*/
import { RichlistComponent } from './richlist.component';

describe('RichlistComponent', () => {
  let component: RichlistComponent;
  let fixture: ComponentFixture<RichlistComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
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

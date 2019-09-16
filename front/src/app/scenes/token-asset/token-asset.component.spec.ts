import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {TokenAssetComponent} from './token-asset.component';
import {AppModule} from '../../app.module';

describe('TokenAssetComponent', () => {
  let component: TokenAssetComponent;
  let fixture: ComponentFixture<TokenAssetComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [AppModule]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TokenAssetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});

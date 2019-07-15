import {inject, TestBed} from '@angular/core/testing';

import {WalletGuard} from './wallet.guard';
import {AppModule} from '../app.module';

describe('WalletGuard', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AppModule],
    });
  });

  it('should ...', inject([WalletGuard], (guard: WalletGuard) => {
    expect(guard).toBeTruthy();
  }));
});

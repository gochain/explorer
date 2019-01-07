import {TestBed} from '@angular/core/testing';
import {WalletService} from './wallet.service';
import {WalletModule} from './wallet.module';
import {AppModule} from '../../app.module';

describe('WalletService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    imports: [AppModule, WalletModule],
  }));

  it('should be created', () => {
    const service: WalletService = TestBed.get(WalletService);
    expect(service).toBeTruthy();
  });
});

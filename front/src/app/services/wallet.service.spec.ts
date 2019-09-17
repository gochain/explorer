import {TestBed} from '@angular/core/testing';
import {WalletService} from './wallet.service';
import {AppModule} from '../app.module';

describe('WalletService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    imports: [AppModule],
  }));

  it('should be created', () => {
    const service: WalletService = TestBed.get(WalletService);
    expect(service).toBeTruthy();
  });
});

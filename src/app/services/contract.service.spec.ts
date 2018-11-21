import { TestBed, inject } from '@angular/core/testing';

import { ContractService } from './contract.service';
import { AppModule } from '../app.module';

describe('ContractService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [AppModule],
      providers: [ContractService],
    });
  });

  it('should be created', inject([ContractService], (service: ContractService) => {
    expect(service).toBeTruthy();
  }));
});

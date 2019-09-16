import { TestBed } from '@angular/core/testing';

import { MetaService } from './meta.service';

describe('MetaService', () => {
  beforeEach(() => TestBed.configureTestingModule({}));

  it('should be created', () => {
    const service: MetaService = TestBed.get(MetaService);
    expect(service).toBeTruthy();
  });
});

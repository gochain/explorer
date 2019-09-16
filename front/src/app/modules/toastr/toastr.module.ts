import { ModuleWithProviders, NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ToastrComponent } from './toastr.component';
import { ToastrService } from './toastr.service';

@NgModule({
  declarations: [ToastrComponent],
  imports: [CommonModule],
  exports: [ToastrComponent],
})
export class ToastrModule {
  static forRoot(): ModuleWithProviders {
    return {
      ngModule: ToastrModule,
      providers: [ToastrService]
    };
  }

  static forChild(): ModuleWithProviders {
    return {
      ngModule: ToastrModule,
      providers: []
    };
  }
}

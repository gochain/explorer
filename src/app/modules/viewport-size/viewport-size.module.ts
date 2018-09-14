import { ModuleWithProviders, NgModule } from '@angular/core';
import { ViewportSizeDirective } from './viewport-size.directive';
import { ViewportSizeService } from './viewport-size.service';
import { CommonModule } from '@angular/common';
import { IConfig } from './config.interface';

@NgModule({
  declarations: [ViewportSizeDirective],
  imports: [CommonModule],
  exports: [ViewportSizeDirective],
  providers: [ViewportSizeService]
})
export class ViewportSizeModule {
  static forRoot(config: IConfig): ModuleWithProviders {
    return {
      ngModule: ViewportSizeModule,
      providers: [ViewportSizeModule, {provide: 'config', useValue: config}]
    };
  }
}

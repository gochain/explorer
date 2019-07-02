/*MODULES*/
import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
/*SERVICES*/
import {CommonService} from '../../services/common.service';
import {ApiService} from '../../services/api.service';
import {LayoutService} from '../../services/layout.service';
import {MetaService} from '../../services/meta.service';
/*COMPONENTS*/
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {HttpClientModule} from '@angular/common/http';
import {PipesModule} from '../pipes/pipes.module';
import {DirectiveModule} from '../../directives/directives.module';
import {TabsModule} from '../tabs/tabs.module';
import {SliderModule} from '../slider/slider.module';
import {ViewportSizeModule} from '../viewport-size/viewport-size.module';
import {VIEWPORT_SIZES} from '../viewport-size/contants';
import {NgProgressModule} from '@ngx-progressbar/core';
import {NgProgressHttpModule} from '@ngx-progressbar/http';
import {ToastrModule} from '../toastr/toastr.module';

@NgModule({
  declarations: [
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    PipesModule,
    DirectiveModule,
    TabsModule,
    SliderModule,
    ViewportSizeModule.forRoot(VIEWPORT_SIZES),
    NgProgressModule.withConfig({
      trickleSpeed: 200,
      min: 20,
      meteor: false,
      spinner: false
    }),
    NgProgressHttpModule,
    ToastrModule.forRoot(),
  ],
  exports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    PipesModule,
    DirectiveModule,
    TabsModule,
    SliderModule,
    ViewportSizeModule,
    NgProgressModule,
    NgProgressHttpModule,
    ToastrModule,
  ],
})
export class SharedModule {
}

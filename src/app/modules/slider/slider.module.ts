import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SliderComponent } from './slider.component';
import {PipesModule} from '../pipes.module';

@NgModule({
  declarations: [SliderComponent],
  imports: [CommonModule, PipesModule],
  exports: [SliderComponent]
})
export class SliderModule { }

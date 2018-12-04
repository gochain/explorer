import {NgModule} from '@angular/core';
import {InfinityScrollDirective} from './infinity-scroll.directive';
import {RecaptchaDirective} from './recaptcha.directive';

@NgModule({
  declarations: [
    InfinityScrollDirective,
    RecaptchaDirective
  ],
  imports: [],
  exports: [
    InfinityScrollDirective,
    RecaptchaDirective
  ]
})
export class DirectiveModule {
}

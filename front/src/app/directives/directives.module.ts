import {NgModule} from '@angular/core';
import {InfinityScrollDirective} from './infinity-scroll.directive';
import {ReCaptchaDirective} from './recaptcha.directive';

@NgModule({
  declarations: [
    InfinityScrollDirective,
    ReCaptchaDirective
  ],
  imports: [],
  exports: [
    InfinityScrollDirective,
    ReCaptchaDirective
  ]
})
export class DirectiveModule {
}

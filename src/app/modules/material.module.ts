import {
  MatProgressBarModule,
  MatSelectModule,
  MatSnackBarModule,
  MatTabsModule
} from '@angular/material';
import {NgModule} from '@angular/core';

@NgModule({
  imports: [
    MatProgressBarModule,
    MatSnackBarModule,
    MatSelectModule,
    MatTabsModule
  ],
  exports: [
    MatProgressBarModule,
    MatSnackBarModule,
    MatSelectModule,
    MatTabsModule
  ],
})
export class MaterialModule {
}
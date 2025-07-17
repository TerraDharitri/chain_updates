import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-label',
  standalone: true,
  imports: [CommonModule],
  template: `<span class="font-semibold"><ng-content></ng-content></span>`
})
export class LabelComponent {
} 
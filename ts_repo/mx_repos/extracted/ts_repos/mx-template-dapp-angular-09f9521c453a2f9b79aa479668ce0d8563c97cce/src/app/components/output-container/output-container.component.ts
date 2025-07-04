import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-output-container',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './output-container.component.html',
  styleUrls: ['./output-container.component.css']
})
export class OutputContainerComponent {
  @Input() isLoading: boolean = false;
  @Input() className: string = 'p-4';
  @Input() dataTestId?: string;
} 
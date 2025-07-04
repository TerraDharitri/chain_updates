import { Component, Input } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CardComponent } from '../card/card.component';
import { WidgetType } from '../../../types/widget.types';

@Component({
  selector: 'app-widget',
  standalone: true,
  imports: [CommonModule, CardComponent],
  templateUrl: './widget.component.html',
  styleUrls: ['./widget.component.css']
})
export class WidgetComponent {
  @Input() widgetConfig!: WidgetType;
} 
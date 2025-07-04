import { Component, Input, OnInit, OnDestroy, OnChanges } from '@angular/core';
import { CommonModule } from '@angular/common';

export interface FormatAmountProps {
  value: string;
  className?: string;
  'data-testid'?: string;
  showLabel?: boolean;
  rewaLabel?: string;
  digits?: number;
  decimals?: number;
}

@Component({
  selector: 'app-format-amount',
  standalone: true,
  imports: [CommonModule],
  template: `
    <span [class]="className" [attr.data-testid]="dataTestId">
      <span class="amount-integer">{{ valueInteger }}</span
      ><span class="amount-decimal" *ngIf="valueDecimal"
        >.{{ valueDecimal }}</span
      >
      <span class="amount-label" *ngIf="showLabel"> {{ label }}</span>
    </span>
  `,
  styleUrls: ['./format-amount.component.css'],
})
export class FormatAmountComponent implements OnInit, OnDestroy, OnChanges {
  @Input() value: string = '';
  @Input() className?: string;
  @Input() dataTestId?: string;
  @Input() showLabel?: boolean = true;
  @Input() rewaLabel?: string = 'REWA';
  @Input() digits?: number = 4;
  @Input() decimals?: number = 18;

  // Computed properties for display
  isValid: boolean = true;
  valueDecimal: string = '';
  valueInteger: string = '';
  label: string = '';

  ngOnInit() {
    this.updateFormattedData();
  }

  ngOnDestroy() {
    // Cleanup if needed
  }

  ngOnChanges() {
    this.updateFormattedData();
  }

  private updateFormattedData() {
    try {
      // Simple amount formatting without external dependencies
      const amount = this.value || '0';

      // Validate if it's a valid number
      const numericValue = parseFloat(amount);
      this.isValid = !isNaN(numericValue) && isFinite(numericValue);

      if (this.isValid && numericValue > 0) {
        // Format the amount: divide by 10^decimals to get the actual value
        const actualValue = numericValue / Math.pow(10, this.decimals || 18);

        // Split into integer and decimal parts
        const formatted = actualValue.toFixed(this.digits || 4);
        const parts = formatted.split('.');

        this.valueInteger = parts[0] || '0';
        this.valueDecimal = parts[1] || '';
        this.label = this.rewaLabel || 'REWA';
      } else {
        this.valueInteger = '0';
        this.valueDecimal = '';
        this.label = this.rewaLabel || 'REWA';
      }
    } catch (error) {
      console.error('Error formatting amount:', error);
      // Fallback values
      this.isValid = false;
      this.valueDecimal = '';
      this.valueInteger = this.value || '0';
      this.label = this.rewaLabel || 'REWA';
    }
  }
}

import { Component, Input } from '@angular/core';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-mx-link',
  standalone: true,
  imports: [RouterModule],
  templateUrl: './mx-link.component.html',
  styleUrls: ['./mx-link.component.css']
})
export class MxLinkComponent {
  @Input() to!: string;
  @Input() className = 'inline-block rounded-lg px-3 py-2 text-center hover:no-underline my-0 bg-blue-600 text-white hover:bg-blue-700 ml-2 mr-0';
} 
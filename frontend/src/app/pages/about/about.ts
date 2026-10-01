import { Component } from '@angular/core';
import { amenities } from '../../core/data/hotel';
import { PageHeader } from '../../layout/page-header/page-header';

@Component({
  selector: 'app-about',
  imports: [PageHeader],
  templateUrl: './about.html',
  styleUrl: './about.css',
})
export class About {
  readonly amenities = amenities;
}

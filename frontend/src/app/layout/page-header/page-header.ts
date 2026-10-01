import { Component } from '@angular/core';
import { hotel } from '../../core/data/hotel';

@Component({
  selector: 'app-page-header',
  templateUrl: './page-header.html',
  styleUrl: './page-header.css',
})
export class PageHeader {
  readonly hotel = hotel;
}

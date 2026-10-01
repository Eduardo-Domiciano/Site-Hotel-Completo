import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';
import { hotel } from '../../core/data/hotel';

@Component({
  selector: 'app-site-footer',
  imports: [RouterLink],
  templateUrl: './site-footer.html',
  styleUrl: './site-footer.css',
})
export class SiteFooter {
  readonly hotel = hotel;
}

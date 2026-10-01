import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { careerAreas } from '../../core/data/hotel';
import { ApiService } from '../../core/services/api';
import { PageHeader } from '../../layout/page-header/page-header';

@Component({
  selector: 'app-career',
  imports: [PageHeader, FormsModule],
  templateUrl: './career.html',
  styleUrl: './career.css',
})
export class Career {
  private readonly api = inject(ApiService);
  readonly areas = careerAreas;
  form = { area: '', name: '', phone: '', email: '', fileName: '' };
  private file: File | null = null;
  message = '';

  onFile(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.file = input.files?.[0] ?? null;
    this.form.fileName = this.file?.name ?? '';
  }

  submit(): void {
    if (!this.form.area || !this.form.name || !this.form.email) {
      this.message = 'Preencha área, nome e e-mail.';
      return;
    }
    const data = new FormData();
    data.append('area', this.form.area);
    data.append('name', this.form.name);
    data.append('email', this.form.email);
    if (this.form.phone) data.append('phone', this.form.phone);
    if (this.file) data.append('resume', this.file);

    this.api.createJobApplication(data).subscribe({
      next: () => {
        this.message = `Candidatura de ${this.form.name} para ${this.form.area} enviada.`;
      },
      error: (err) => {
        this.message = err?.error?.error?.message ?? 'Falha ao enviar candidatura.';
      },
    });
  }
}

import {
  Component,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

export type Prioridad = 'alta' | 'media' | 'baja';
export type Filtro    = 'todas' | 'pendientes' | 'completadas';

export interface Tarea {
  id: number;
  texto: string;
  prioridad: Prioridad;
  completada: boolean;
  editando: boolean;
  textoEdicion: string;
}

@Component({
  selector: 'app-tareas',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './tareas.component.html',
  styleUrls: ['./tareas.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class TareasComponent {

  private nextId = 4;

  tareas = signal<Tarea[]>([
    this.nuevaTarea('Escribir reseña de la RTX 5090',              'alta'),
    this.nuevaTarea('Revisar borrador de Gloria sobre indies',      'media'),
    this.nuevaTarea('Subir imágenes de portada a Supabase Storage', 'baja'),
  ]);

  filtro       = signal<Filtro>('todas');
  errorMsg     = signal<string>('');
  nuevoTexto   = '';
  nuevaPrioridad: Prioridad = 'media';

  total       = computed(() => this.tareas().length);
  completadas = computed(() => this.tareas().filter(t => t.completada).length);
  pendientes  = computed(() => this.total() - this.completadas());
  progreso    = computed(() =>
    this.total() ? Math.round((this.completadas() / this.total()) * 100) : 0
  );

  tareasFiltradas = computed(() => {
    const f = this.filtro();
    return this.tareas().filter(t =>
      f === 'todas'       ? true :
      f === 'pendientes'  ? !t.completada :
                             t.completada
    );
  });

  private validar(texto: string): string {
    const t = texto.trim();
    if (!t)             return 'La tarea no puede estar vacía.';
    if (t.length < 3)   return 'Mínimo 3 caracteres.';
    if (t.length > 80)  return 'Máximo 80 caracteres.';
    const existe = this.tareas().some(
      ta => ta.texto.toLowerCase() === t.toLowerCase()
    );
    if (existe) return 'Esa tarea ya existe.';
    return '';
  }

  agregar(): void {
    const error = this.validar(this.nuevoTexto);
    if (error) { this.errorMsg.set(error); return; }
    this.errorMsg.set('');
    this.tareas.update(arr => [
      ...arr,
      this.nuevaTarea(this.nuevoTexto.trim(), this.nuevaPrioridad),
    ]);
    this.nuevoTexto = '';
  }

  toggleCompletar(id: number): void {
    this.tareas.update(arr =>
      arr.map(t => t.id === id ? { ...t, completada: !t.completada } : t)
    );
  }

  eliminar(id: number): void {
    this.tareas.update(arr => arr.filter(t => t.id !== id));
  }

  limpiarCompletadas(): void {
    this.tareas.update(arr => arr.filter(t => !t.completada));
  }

  iniciarEdicion(tarea: Tarea): void {
    this.tareas.update(arr =>
      arr.map(t => ({
        ...t,
        editando:     t.id === tarea.id,
        textoEdicion: t.id === tarea.id ? t.texto : t.textoEdicion,
      }))
    );
  }

  guardarEdicion(tarea: Tarea): void {
    const nuevo = tarea.textoEdicion.trim();
    this.tareas.update(arr =>
      arr.map(t =>
        t.id === tarea.id
          ? { ...t, texto: nuevo || t.texto, editando: false }
          : t
      )
    );
  }

  cancelarEdicion(id: number): void {
    this.tareas.update(arr =>
      arr.map(t => t.id === id ? { ...t, editando: false } : t)
    );
  }

  onEditKeydown(event: KeyboardEvent, tarea: Tarea): void {
    if (event.key === 'Enter')  this.guardarEdicion(tarea);
    if (event.key === 'Escape') this.cancelarEdicion(tarea.id);
  }

  setFiltro(f: Filtro): void { this.filtro.set(f); }
  limpiarError(): void { if (this.errorMsg()) this.errorMsg.set(''); }
  onEnter(event: KeyboardEvent): void { if (event.key === 'Enter') this.agregar(); }
  trackById(_: number, t: Tarea): number { return t.id; }

  private nuevaTarea(texto: string, prioridad: Prioridad): Tarea {
    return {
      id: this.nextId++,
      texto,
      prioridad,
      completada: false,
      editando: false,
      textoEdicion: '',
    };
  }
}

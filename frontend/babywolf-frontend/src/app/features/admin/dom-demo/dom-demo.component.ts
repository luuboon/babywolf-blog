import {
  Component,
  signal,
  computed,
  ChangeDetectionStrategy,
} from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

export interface PostCard {
  id: number;
  titulo: string;
  categoria: string;
  destacada: boolean;
}

export interface LogEntry {
  id: number;
  metodo: string;
  detalle: string;
  ts: string;
}

@Component({
  selector: 'app-dom-demo',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './dom-demo.component.html',
  styleUrls: ['./dom-demo.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DomDemoComponent {

  private nextId = 4;
  private logId  = 1;

  tituloBlog  = signal('BabyWolf Blog');
  temaOscuro  = signal(false);
  colorAcento = signal('#e94560');
  tamFuente   = signal(16);
  bannerVisible = signal(false);
  bannerTexto   = signal('');
  cardEditandoId = signal<number | null>(null);
  cardEditTexto  = '';
  nuevoTitulo    = '';
  nuevaCategoria = 'gaming';

  cards = signal<PostCard[]>([
    { id: 1, titulo: 'Análisis: RTX 5090',     categoria: 'hardware',   destacada: false },
    { id: 2, titulo: 'Top 10 juegos indie',     categoria: 'gaming',     destacada: false },
    { id: 3, titulo: 'Angular Signals vs RxJS', categoria: 'tutoriales', destacada: false },
  ]);

  log = signal<LogEntry[]>([]);

  totalCards = computed(() => this.cards().length);

  actualizarTitulo(valor: string): void {
    this.tituloBlog.set(valor || 'BabyWolf Blog');
    this.registrar('textContent', `tituloBlog = "${this.tituloBlog()}"`);
  }

  toggleTema(): void {
    this.temaOscuro.update(v => !v);
    this.registrar('classList.toggle("dark")', `tema → ${this.temaOscuro() ? 'oscuro' : 'claro'}`);
  }

  cambiarAcento(color: string): void {
    this.colorAcento.set(color);
    this.registrar('style.setProperty("--accent")', `color → ${color}`);
  }

  cambiarFuente(tam: number): void {
    this.tamFuente.set(tam);
    this.registrar('style.fontSize', `${tam}px`);
  }

  agregarCard(): void {
    if (!this.nuevoTitulo.trim()) return;
    const card: PostCard = { id: this.nextId++, titulo: this.nuevoTitulo.trim(), categoria: this.nuevaCategoria, destacada: false };
    this.cards.update(arr => [...arr, card]);
    this.registrar('createElement + appendChild', `card "${card.titulo}" → #postsGrid`);
    this.nuevoTitulo = '';
  }

  onEnter(e: KeyboardEvent): void { if (e.key === 'Enter') this.agregarCard(); }

  eliminarCard(id: number): void {
    const card = this.cards().find(c => c.id === id);
    this.cards.update(arr => arr.filter(c => c.id !== id));
    this.registrar('element.remove()', `card "${card?.titulo}" eliminada`);
  }

  duplicarCard(card: PostCard): void {
    const clon: PostCard = { ...card, id: this.nextId++, titulo: card.titulo + ' (copia)' };
    this.cards.update(arr => {
      const idx = arr.findIndex(c => c.id === card.id);
      const next = [...arr];
      next.splice(idx + 1, 0, clon);
      return next;
    });
    this.registrar('cloneNode(true) + insertBefore', `"${card.titulo}" duplicada`);
  }

  toggleDestacada(id: number): void {
    this.cards.update(arr => arr.map(c => c.id === id ? { ...c, destacada: !c.destacada } : c));
    const c = this.cards().find(c => c.id === id);
    this.registrar('classList.toggle("destacada")', `"${c?.titulo}" → destacada: ${c?.destacada}`);
  }

  iniciarEdicion(card: PostCard): void {
    this.cardEditandoId.set(card.id);
    this.cardEditTexto = card.titulo;
    this.registrar('replaceWith(input)', `edición in-place de "${card.titulo}"`);
  }

  guardarEdicion(id: number): void {
    const nuevo = this.cardEditTexto.trim();
    if (nuevo) {
      this.cards.update(arr => arr.map(c => c.id === id ? { ...c, titulo: nuevo } : c));
      this.registrar('input.replaceWith(span)', `título guardado: "${nuevo}"`);
    }
    this.cardEditandoId.set(null);
  }

  cancelarEdicion(): void { this.cardEditandoId.set(null); }

  onEditKeydown(e: KeyboardEvent, id: number): void {
    if (e.key === 'Enter')  this.guardarEdicion(id);
    if (e.key === 'Escape') this.cancelarEdicion();
  }

  mostrarBanner(): void {
    this.bannerTexto.set('🎉 Banner creado dinámicamente con prepend() / signal');
    this.bannerVisible.set(true);
    this.registrar('prepend()', 'banner insertado al inicio del contenido');
    setTimeout(() => this.cerrarBanner(), 4000);
  }

  cerrarBanner(): void {
    this.bannerVisible.set(false);
    this.registrar('element.remove()', 'banner eliminado');
  }

  limpiarLog(): void { this.log.set([]); }

  trackById = (_: number, item: { id: number }) => item.id;

  private registrar(metodo: string, detalle: string): void {
    const entry: LogEntry = { id: this.logId++, metodo, detalle, ts: new Date().toLocaleTimeString('es-MX') };
    this.log.update(arr => [entry, ...arr].slice(0, 20));
  }
}

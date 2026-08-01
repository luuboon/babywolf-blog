import { Injectable } from '@angular/core';

interface SeasonalTheme {
  className: string;
  label: string;
  /** Mes (1-12) de inicio y fin, inclusive. Cruza fin de año en diciembre-enero. */
  months: number[];
}

/**
 * Practica 7-8: evento calendarizado que cambia el fondo/banner del sitio
 * según la fecha actual (sin intervención del usuario).
 */
const THEMES: SeasonalTheme[] = [
  { className: 'season-halloween', label: 'Halloween 🎃', months: [10] },
  { className: 'season-navidad', label: 'Navidad 🎄', months: [12] },
  { className: 'season-verano', label: 'Verano ☀️', months: [6, 7, 8] },
];
const DEFAULT_THEME = 'season-default';

@Injectable({ providedIn: 'root' })
export class SeasonalThemeService {
  apply(date: Date = new Date()): string {
    if (typeof document === 'undefined') return DEFAULT_THEME;

    const month = date.getMonth() + 1;
    const theme = THEMES.find(t => t.months.includes(month));
    const className = theme?.className ?? DEFAULT_THEME;

    document.body.classList.remove(...THEMES.map(t => t.className), DEFAULT_THEME);
    document.body.classList.add(className);
    return className;
  }
}

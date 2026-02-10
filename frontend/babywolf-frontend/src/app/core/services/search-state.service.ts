import { Injectable, signal, computed } from '@angular/core';

@Injectable({
    providedIn: 'root'
})
export class SearchStateService {
    // Private writable signal
    private searchSignal = signal<string>('');

    // Public readonly signal for consumption
    readonly searchQuery = computed(() => this.searchSignal());

    setSearchQuery(query: string): void {
        this.searchSignal.set(query);
    }

    clearSearch(): void {
        this.searchSignal.set('');
    }
}

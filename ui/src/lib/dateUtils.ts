import type { PriceEntry } from '$lib/pricesApi';

export function formatDateDay(d: Date): string {
	return d.toLocaleDateString('fi-FI');
}

export function formatDateTime(d: Date): string {
	const time = d.toLocaleTimeString('fi-FI');
	return time.replaceAll('.00', '');
}

export function isNow(p: PriceEntry) {
	const now = new Date();
	return p.s < now && p.e > now;
}

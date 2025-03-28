import type { PriceEntry } from '$lib/pricesApi';

export function formatDateDay(d: Date): string {
	return d.toLocaleDateString('fi-FI');
}

export function formatDateTime(d: Date): string {
	const time = d.toLocaleTimeString('fi-FI');
	return time.replaceAll('.00', '');
}

export function isNow(p: PriceEntry, now?: Date) {
	if (!now) {
		now = new Date();
	}
	return p.s < now && p.e > now;
}

export function formatMillis(m: number): string {
	return new Date(m).toISOString().slice(11, 19);
}

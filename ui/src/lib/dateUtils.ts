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

function pad(n: number) {
	if (n < 10) {
		return '0' + n.toString();
	}
	return n.toString();
}
export function formatMillis(ms: number): string {
	let seconds = ms / 1000;
	const hours = Math.floor(seconds / 3600);
	seconds = seconds % 3600;
	const minutes = Math.floor(seconds / 60);
	seconds = seconds % 60;

	return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}

function minutesToMillis(minutes: number) {
	return 1000 * 60 * minutes;
}
export const millisBeforeToTryFetch = minutesToMillis(20);

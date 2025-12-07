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
		return '0' + n.toString(10);
	}
	return n.toString(10);
}
export function formatSeconds(seconds: number): string {
	const hours = Math.floor(seconds / 3600);
	seconds = seconds % 3600;
	const minutes = Math.floor(seconds / 60);
	seconds = Math.floor(seconds % 60);

	return `${pad(hours)}:${pad(minutes)}:${pad(seconds)}`;
}

function minutesToSeconds(minutes: number) {
	return 60 * minutes;
}
export const secondsBeforeToTryFetch = minutesToSeconds(20);

export function isNextDayVisible(prices: PriceEntry[], now: Date): boolean {
	return prices.some((p) => {
		const s = p.s;
		// detect if `s` is a later calendar day than `now`
		const isLaterDay =
			s.getFullYear() > now.getFullYear() ||
			(s.getFullYear() === now.getFullYear() && s.getMonth() > now.getMonth()) ||
			(s.getFullYear() === now.getFullYear() &&
				s.getMonth() === now.getMonth() &&
				s.getDate() > now.getDate());

		// only consider later-day prices that are at or after 03:00 local time
		return isLaterDay && s.getHours() >= 3;
	});
}

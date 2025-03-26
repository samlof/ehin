export function formatDateDay(d: Date): string {
	return d.toLocaleDateString('fi-FI');
}

export function formatDateTime(d: Date): string {
	const time = d.toLocaleTimeString('fi-FI');
	return time.replaceAll('.00', '');
}

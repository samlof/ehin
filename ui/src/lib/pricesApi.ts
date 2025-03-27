import { PUBLIC_API_URL } from '$env/static/public';

export interface ResponsePriceEntry {
	p: number;
	s: string;
	e: string;
}

export interface PriceEntry {
	p: number;
	s: Date;
	e: Date;
}

type fetchType = {
	(input: RequestInfo | URL, init?: RequestInit): Promise<Response>;
	(input: RequestInfo | URL, init?: RequestInit<RequestInitCfProperties>): Promise<Response>;
};

const url = PUBLIC_API_URL + '/api/prices';
export async function getPrices(fetch?: fetchType): Promise<PriceEntry[]> {
	if (!fetch) {
		if (!window) {
			console.error('No fetch param but also no window');
		}
		fetch = window.fetch;
	}

	const res = await fetch(url + '/' + getDateForApi());
	const prices: ResponsePriceEntry[] = await res.json();

	return prices.map((p) => ({
		p: p.p,
		e: new Date(p.e),
		s: new Date(p.s),
	}));
}

function getDateForApi() {
	const today = new Date();
	const offset = today.getTimezoneOffset();
	const today2 = new Date(today.getTime() - offset * 60 * 1000);
	const todayDate = today2.toISOString().split('T')[0];
	return todayDate;
}

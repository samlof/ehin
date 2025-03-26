import { API_URL } from '$env/static/private';

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

const url = API_URL + '/api/prices/';
export async function getPrices(fetch?: fetchType): Promise<PriceEntry[]> {
	if (!fetch) {
		if (!window) {
			console.error('No fetch param but also no window');
		}
		fetch = window.fetch;
	}
	const res = await fetch(url + `/2025-03-26`);
	const prices: ResponsePriceEntry[] = await res.json();

	return prices.map((p) => ({
		p: p.p,
		e: new Date(p.e),
		s: new Date(p.s)
	}));
}

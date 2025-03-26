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
	const res = await fetch(url + `/2025-03-26`);
	const prices: ResponsePriceEntry[] = await res.json();

	return prices.map((p) => ({
		p: p.p,
		e: new Date(p.e),
		s: new Date(p.s)
	}));
}

export function demoSeries() {
	return [
		{
			date: new Date('2025-02-24T22:00:00.000Z'),
			value: 44,
			baseline: 100
		},
		{
			date: new Date('2025-02-25T22:00:00.000Z'),
			value: 57,
			baseline: 46
		},
		{
			date: new Date('2025-02-26T22:00:00.000Z'),
			value: 72,
			baseline: 46
		},
		{
			date: new Date('2025-02-27T22:00:00.000Z'),
			value: 35,
			baseline: 46
		},
		{
			date: new Date('2025-02-28T22:00:00.000Z'),
			value: 25,
			baseline: 20
		},
		{
			date: new Date('2025-03-01T22:00:00.000Z'),
			value: 97,
			baseline: 28
		},
		{
			date: new Date('2025-03-02T22:00:00.000Z'),
			value: 95,
			baseline: 81
		},
		{
			date: new Date('2025-03-03T22:00:00.000Z'),
			value: 58,
			baseline: 48
		},
		{
			date: new Date('2025-03-04T22:00:00.000Z'),
			value: 33,
			baseline: 68
		},
		{
			date: new Date('2025-03-05T22:00:00.000Z'),
			value: 41,
			baseline: 55
		},
		{
			date: new Date('2025-03-06T22:00:00.000Z'),
			value: 37,
			baseline: 32
		},
		{
			date: new Date('2025-03-07T22:00:00.000Z'),
			value: 58,
			baseline: 48
		},
		{
			date: new Date('2025-03-08T22:00:00.000Z'),
			value: 48,
			baseline: 89
		},
		{
			date: new Date('2025-03-09T22:00:00.000Z'),
			value: 48,
			baseline: 33
		},
		{
			date: new Date('2025-03-10T22:00:00.000Z'),
			value: 73,
			baseline: 36
		},
		{
			date: new Date('2025-03-11T22:00:00.000Z'),
			value: 99,
			baseline: 29
		},
		{
			date: new Date('2025-03-12T22:00:00.000Z'),
			value: 40,
			baseline: 60
		},
		{
			date: new Date('2025-03-13T22:00:00.000Z'),
			value: 21,
			baseline: 85
		},
		{
			date: new Date('2025-03-14T22:00:00.000Z'),
			value: 85,
			baseline: 87
		},
		{
			date: new Date('2025-03-15T22:00:00.000Z'),
			value: 77,
			baseline: 37
		},
		{
			date: new Date('2025-03-16T22:00:00.000Z'),
			value: 22,
			baseline: 93
		},
		{
			date: new Date('2025-03-17T22:00:00.000Z'),
			value: 67,
			baseline: 48
		},
		{
			date: new Date('2025-03-18T22:00:00.000Z'),
			value: 86,
			baseline: 75
		},
		{
			date: new Date('2025-03-19T22:00:00.000Z'),
			value: 97,
			baseline: 21
		},
		{
			date: new Date('2025-03-20T22:00:00.000Z'),
			value: 79,
			baseline: 55
		},
		{
			date: new Date('2025-03-21T22:00:00.000Z'),
			value: 71,
			baseline: 91
		},
		{
			date: new Date('2025-03-22T22:00:00.000Z'),
			value: 77,
			baseline: 20
		},
		{
			date: new Date('2025-03-23T22:00:00.000Z'),
			value: 71,
			baseline: 74
		},
		{
			date: new Date('2025-03-24T22:00:00.000Z'),
			value: 21,
			baseline: 46
		},
		{
			date: new Date('2025-03-25T22:00:00.000Z'),
			value: 77,
			baseline: 38
		}
	];
}

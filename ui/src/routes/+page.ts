import { getPrices } from '$lib/pricesApi';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	return {
		prices: await getPrices(fetch)
	};
};

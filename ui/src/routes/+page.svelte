<script lang="ts">
	import { getPrices } from '$lib/pricesApi';
	import PriceChart from '$lib/PriceChart.svelte';

	let pricesPromise = getPrices();

	function onfocus() {
		console.log('update chart');
		const prom = getPrices().then((p) => {
			pricesPromise = prom;
			return p;
		});
	}
</script>

<svelte:window {onfocus} />

<article class="prose lg:prose-xl">
	<h1>Sähkö hinnat</h1>
	{#await pricesPromise}
		Loading...
	{:then prices}
		<PriceChart {prices} />
	{:catch error}
		Failed to load prices {error.message}
	{/await}
</article>

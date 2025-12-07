<script lang="ts">
	import { getPrices } from '$lib/pricesApi';
	import PriceChart from '$lib/PriceChart.svelte';
	import Settings from '$lib/Settings.svelte';

	let pricesPromise = getPrices();

	let lastUpdate = new Date();

	function updatePrices() {
		const oneMinute = 1000 * 60;
		if (new Date().getTime() - lastUpdate.getTime() < oneMinute) {
			return;
		}
		lastUpdate = new Date();
		console.log('updating prices');
		const prom = getPrices().then((p) => {
			pricesPromise = prom;
			return p;
		});
	}
	function onfocus() {
		updatePrices();
	}
</script>

<svelte:window {onfocus} />

<article class="main-content">
	<article class="prose lg:prose-xl"><h1>Sähkö hinnat</h1></article>
	{#await pricesPromise}
		Loading...
	{:then prices}
		<PriceChart {prices} {updatePrices} />
	{:catch error}
		Failed to load prices {error.message}
	{/await}

	<Settings />
</article>

<style>
	.main-content {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		width: 100vw !important;
	}
</style>

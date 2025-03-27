<script lang="ts">
	import { getPrices } from '$lib/pricesApi';
	import PriceChart from '$lib/PriceChart.svelte';

	const pricesPromise = getPrices();
</script>

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

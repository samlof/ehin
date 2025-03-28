<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { isNow } from '$lib/dateUtils';
	import { breakpoint, type breakpointVals } from '$lib/mediaQuery.svelte';
	import { chartConfig, setupChart } from '$lib/PriceChartCanvasSetup.svelte';
	import type { PriceEntry } from '$lib/pricesApi';

	interface Props {
		prices: PriceEntry[];
	}

	let { prices }: Props = $props();

	const sizes: { [key in breakpointVals]: number } = {
		xs: 20,
		s: 15,
		m: 10,
		l: 5,
		xl: 1,
		xxl: 1,
	};
	const filteredPrices = $derived(prices.slice(sizes[breakpoint()]));

	const nextDayVisible = $derived(prices[prices.length - 1].s.getDate() !== new Date().getDate());

	const config = $derived(chartConfig(filteredPrices));
	const priceNow = $derived(filteredPrices.find(isNow));
	setupChart('priceChart', () => config);
</script>

{#if priceNow}
	<article class="prose lg:prose-xl pt-5">
		<h2>Hinta nyt {formatPrice(priceNow.p)} c/kWh</h2>
	</article>
{/if}

<div class="canvas-container"><canvas id="priceChart"></canvas></div>

<style>
	.canvas-container {
		width: min(100vw, 1200px);
	}
</style>

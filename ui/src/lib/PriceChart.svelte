<script lang="ts">
	import { chartConfig, setupChart } from '$lib/PriceChartCanvasSetup.svelte';
	import type { PriceEntry } from '$lib/pricesApi';

	interface Props {
		prices: PriceEntry[];
	}

	let { prices }: Props = $props();
	const config = $derived(chartConfig(prices));
	const { update } = setupChart('priceChart', () => config);

	function onfocus() {
		console.log('update chart');
		update();
	}
</script>

<svelte:window {onfocus} />

<div style="width: 800px;"><canvas id="priceChart"></canvas></div>

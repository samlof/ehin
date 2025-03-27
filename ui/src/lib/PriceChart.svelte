<script lang="ts">
	import { setupChart } from '$lib/PriceChartCanvasSetup';
	import type { PriceEntry } from '$lib/pricesApi';

	interface Props {
		prices: PriceEntry[];
	}

	let { prices }: Props = $props();
	const { update } = setupChart(prices);

	function onfocus() {
		console.log('update chart');
		update();
	}
</script>

<svelte:window {onfocus} />

<div style="width: 800px;"><canvas id="priceChart"></canvas></div>

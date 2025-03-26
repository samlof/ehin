<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { formatDateTime } from '$lib/dateUtils';
	import type { PriceEntry } from '$lib/pricesApi';
	import { Chart } from 'chart.js';
	import { onMount } from 'svelte';

	interface Props {
		prices: PriceEntry[];
	}

	let { prices }: Props = $props();

	onMount(() => {
		new Chart(document.getElementById('acquisitions') as any, {
			type: 'bar',
			data: {
				labels: prices.map((p) => formatDateTime(p.s)),
				datasets: [
					{
						label: 'Acquisitions by year',
						data: prices.map((p) => formatPrice(p.p))
					}
				]
			}
		});
	});
</script>

<div style="width: 800px;"><canvas id="acquisitions"></canvas></div>

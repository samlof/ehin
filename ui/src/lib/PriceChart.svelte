<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { formatDateTime, isNow } from '$lib/dateUtils';
	import type { PriceEntry } from '$lib/pricesApi';
	import {
		Chart,
		Colors,
		BarController,
		CategoryScale,
		LinearScale,
		BarElement,
		Legend,
		Tooltip,
	} from 'chart.js';
	import { onMount } from 'svelte';

	interface Props {
		prices: PriceEntry[];
	}

	let { prices }: Props = $props();

	onMount(() => {
		const defaultColor = 'rgba(54, 162, 235,0.5)';
		const transparentColor = 'rgba(54, 162, 235,0)';
		const blackColor = 'rgba(100, 100, 100,0.5)';
		const todayColor = 'rgb(54, 162, 235)';
		Chart.register(Colors, BarController, BarElement, CategoryScale, LinearScale, Legend, Tooltip);

		const biggest = formatPrice(Math.max(...prices.map((p) => p.p)));

		new Chart(document.getElementById('acquisitions') as any, {
			type: 'bar',
			options: {
				scales: {
					y: {
						beginAtZero: true,
					},
				},
				plugins: {
					legend: {
						labels: {
							filter(item, data) {
								return item.text === 'c/kWh';
							},
						},
					},
					tooltip: {
						filter(e, index, array, data) {
							return e.dataset.label === 'Hover helper' || e.dataset.label === 'c/kWh';
						},
						yAlign: 'center',
						xAlign: 'center',
						displayColors: false,
						callbacks: {
							title(tooltipItems) {
								const p = prices[tooltipItems[0].dataIndex];
								return formatDateTime(p.s) + ' - ' + formatDateTime(p.e);
							},
							label(tooltipItem) {
								return formatPrice(prices[tooltipItem.dataIndex].p) + ' c/kWh';
							},
							beforeLabel(tooltipItem) {
								return '';
							},
							labelPointStyle(tooltipItem) {
								return undefined;
							},
						},
					},
				},
			},
			data: {
				labels: prices.map((p) => formatDateTime(p.s)),
				datasets: [
					{
						label: 'c/kWh',
						data: prices.map((p) => formatPrice(p.p)),
						backgroundColor: prices.map((p) => (isNow(p) ? todayColor : defaultColor)),
						order: 7,
						grouped: false,
					},
					{
						label: 'Nyt',
						data: prices.map((p) => (isNow(p) ? biggest : '0')),
						grouped: false,
						order: 3,
						categoryPercentage: 0.5,
					},
					{
						label: 'Hover helper',
						data: prices.map((p) => biggest),
						backgroundColor: transparentColor,
						order: 20,
						hoverBackgroundColor: defaultColor,
						grouped: false,
					},
					{
						label: 'Day change',
						backgroundColor: prices.map((p) =>
							formatDateTime(p.s) === '0' ? blackColor : transparentColor,
						),
						data: prices.map((p) => (formatDateTime(p.s) === '0' ? biggest : '0')),
						order: 1,
						grouped: false,
						categoryPercentage: 0.1,
					},
				],
			},
		});
	});
</script>

<div style="width: 800px;"><canvas id="acquisitions"></canvas></div>

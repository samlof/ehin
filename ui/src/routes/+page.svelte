<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { formatDateDay, formatDateTime } from '$lib/dateUtils';
	import { demoSeries, getPrices } from '$lib/pricesApi';
	import Chart from 'chart.js/auto';
	import type { PageProps } from './$types';
	import { onMount } from 'svelte';
	import PriceChart from '$lib/PriceChart.svelte';

	let { data }: PageProps = $props();
	const timeNow = new Date();

	const demodata = [
		{ year: 2010, count: 10 },
		{ year: 2011, count: 20 },
		{ year: 2012, count: 15 },
		{ year: 2013, count: 25 },
		{ year: 2014, count: 22 },
		{ year: 2015, count: 30 },
		{ year: 2016, count: 28 }
	];

	const pricesPromise = getPrices();
	const myData = demoSeries();
</script>

<article class="prose lg:prose-xl">
	<h1>Sähkö hinnat</h1>
	{#await pricesPromise}
		Loading...
	{:then prices}
		<PriceChart {prices} />
	{/await}
</article>

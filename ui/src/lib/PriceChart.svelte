<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { formatSeconds, isNow, isNextDayVisible, secondsBeforeToTryFetch } from '$lib/dateUtils';
	import { breakpoint, type breakpointVals } from '$lib/mediaQuery.svelte';
	import { chartConfig, setupChart } from '$lib/PriceChartCanvasSetup.svelte';
	import type { PriceEntry } from '$lib/pricesApi';
	import { onMount } from 'svelte';
	import { SvelteDate } from 'svelte/reactivity';
	import { localSettings } from './stores/localStore.svelte';

	interface Props {
		prices: PriceEntry[];
		updatePrices: () => void;
	}

	let { prices, updatePrices }: Props = $props();

	const sizes: { [key in breakpointVals]: number } = {
		xs: 20,
		s: 15,
		m: 10,
		l: 5,
		xl: 1,
		xxl: 1,
	};
	const filteredPrices = $derived.by(() => {
		if (localSettings.showOnlyAfterNow.value) {
			let nowDate = new Date();
			nowDate.setMinutes(nowDate.getMinutes() - 15);
			const nowTime = nowDate.getTime();
			return prices.filter((p) => p.s.getTime() >= nowTime);
		}
		return prices.slice(sizes[breakpoint()]);
	});

	let now = new SvelteDate();
	const nextDayVisible = $derived(isNextDayVisible(prices, now));

	const secondsUntil = $derived.by(() => {
		/* eslint-disable svelte/prefer-svelte-reactivity */
		const utc12 = new Date();
		utc12.setUTCHours(12, 15, 0, 0);
		if (utc12.getTime() - now.getTime() < 0) {
			utc12.setDate(utc12.getDate() + 1);
		}
		return Math.round((utc12.getTime() - now.getTime()) / 1000);
	});

	onMount(() => {
		const int = setInterval(() => {
			now.setTime(Date.now());
		}, 1000);
		return () => {
			clearInterval(int);
		};
	});

	$effect(() => {
		if (!nextDayVisible && secondsUntil < secondsBeforeToTryFetch) {
			updatePrices();
		}
	});
	const config = $derived(chartConfig(filteredPrices, localSettings.showOnlyAfterNow.value));
	const priceNow = $derived(filteredPrices.find((p) => isNow(p, now)));
	setupChart('priceChart', () => config);
</script>

{#if priceNow}
	<article class="prose lg:prose-xl pt-5">
		<h2>Hinta nyt {formatPrice(priceNow.p)} c/kWh</h2>
	</article>
{/if}

<sub class="py-4">
	{#if nextDayVisible}
		Seuraavat hinnat julkaistaan huomenna noin kello 14
	{:else if secondsUntil < 0}
		Seuraavat hinnat ovat saatavilla hetkenä minä hyvänsä
	{:else}
		Seuraavien hintojen julkaisuun noin {formatSeconds(secondsUntil)}
	{/if}
</sub>

<div class="canvas-container"><canvas id="priceChart"></canvas></div>

<style>
	.canvas-container {
		width: min(100vw, 1200px);
	}
</style>

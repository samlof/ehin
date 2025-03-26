<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { formatDateDay, formatDateTime } from '$lib/dateUtils';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const timeNow = new Date();
</script>

<h1>Sähkö hinnat</h1>
<div>
	{#each data.prices as p}
		<div>
			{#if p.s < timeNow && timeNow < p.e}
				RIGHT NOW
			{/if}
			{formatDateDay(p.s)}
			{formatDateTime(p.s)} - {formatDateTime(p.e)}:
			{formatPrice(p.p)} c/kWh
		</div>
	{/each}
</div>

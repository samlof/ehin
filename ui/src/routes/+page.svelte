<script lang="ts">
	import { formatPrice } from '$lib/calcUtils';
	import { formatDateDay, formatDateTime } from '$lib/dateUtils';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const timeNow = new Date();
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
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

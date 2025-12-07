<script>
	import { localSettings } from './stores/localStore.svelte';
</script>

<div class="settings-root">
	<h3 class="settings-title">Settings</h3>
	<p class="settings-desc">
		Choose which price points to display — either all points, or only those after the current time.
	</p>

	<div class="settings">
		<label>
			<input type="radio" value={false} bind:group={localSettings.showOnlyAfterNow.value} />
			Show All
		</label>

		<label>
			<input type="radio" value={true} bind:group={localSettings.showOnlyAfterNow.value} />
			Show Only After
		</label>
	</div>
</div>

<style>
	.settings {
		display: flex;
		gap: 1rem;
	}

	.settings-title {
		margin: 0 0 0.25rem 0;
		font-size: 1rem;
		font-weight: 600;
	}

	.settings-desc {
		margin: 0 0 0.5rem 0;
		color: var(--muted-color, #666);
		font-size: 0.9rem;
	}

	.settings-root {
		padding-top: 3rem;
	}

	label {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		cursor: pointer;
	}
</style>

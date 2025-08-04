<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';

	let { value = $bindable() }: { value: string[] } = $props();

	const fileTypes = ['exam', 'answer', 'other'];
</script>

<Select.Root type="multiple" name="dateitypen" bind:value>
	<Select.Trigger class="w-[100%] bg-white">
		{#if value.length === 0 || value.length === fileTypes.length}
			Alle Dateiarten
		{:else if value.length > 1}
			Mehere Dateiarten
		{:else}
			{value[0]}
		{/if}
	</Select.Trigger>
	<Select.Content>
		<Select.Group>
			{#each fileTypes as fileType}
				<Select.Item value={fileType} label={fileType}>
					{fileType}
				</Select.Item>
			{/each}
		</Select.Group>
	</Select.Content>
</Select.Root>

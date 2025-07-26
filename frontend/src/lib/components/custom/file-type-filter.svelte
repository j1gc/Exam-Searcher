<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';

	let { value = $bindable() }: { value: string[] } = $props();

	const fileTypes = ['exam', 'answer', 'other'];

	let selectedFileTypes: string[] = $state([]);
</script>

<Select.Root type="multiple" name="dateitypen" bind:value={selectedFileTypes}>
	<Select.Trigger class="w-[180px] bg-white">
		{#if selectedFileTypes.length === 0 || selectedFileTypes.length === fileTypes.length}
			Alle Dateiarten
		{:else if selectedFileTypes.length > 1}
			Mehere Dateiarten
		{:else}
			{selectedFileTypes[0]}
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

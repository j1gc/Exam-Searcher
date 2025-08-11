<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { m } from '$lib/paraglide/messages';

	let { value = $bindable() }: { value: string[] } = $props();

	const fileTypes = ['exam', 'answer', 'other'];
</script>

<Select.Root type="multiple" name="dateitypen" bind:value>
	<Select.Trigger class="w-[100%] bg-white">
		{#if value.length === 0 || value.length === fileTypes.length}
			{m.filter_all()} {m.file_type_text()}
		{:else if value.length > 1}
			{m.filter_multiple()} {m.file_type_text()}
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

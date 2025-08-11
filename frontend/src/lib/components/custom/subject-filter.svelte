<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { m } from '$lib/paraglide/messages';
	import { subjectMapping } from '$lib/subject_mapping.svelte';

	let { value = $bindable() }: { value: number[] } = $props();

	let subjectIdStrings = $state(['-1']);

	$effect(() => {
		if (subjectIdStrings[0] === '-1' && subjectIdStrings.length > 1) {
			subjectIdStrings.shift();
		} else if (subjectIdStrings.length === 0) {
			subjectIdStrings.push('-1');
		}

		value = subjectIdStrings
			.filter((str) => str !== '-1')
			.map((str) => {
				return parseInt(str);
			});
	});
</script>

<Select.Root type="multiple" name="fächer" bind:value={subjectIdStrings}>
	<Select.Trigger class="w-[100%] bg-white">
		{#if (subjectIdStrings[0] === '-1' && subjectIdStrings.length === 1) || subjectIdStrings.length === subjectMapping.size}
			{m.filter_all()} {m.filter_subject_text()}
		{:else if subjectIdStrings.length > 1}
			{m.filter_multiple()} {m.filter_subject_text()}
		{:else}
			{subjectMapping.get(parseInt(subjectIdStrings[0]))}
		{/if}
	</Select.Trigger>
	<Select.Content>
		<Select.Group>
			{#each subjectMapping as subject}
				<Select.Item value={subject[0].toString()} label={subject[1]}>
					{subject[1]}
				</Select.Item>
			{/each}
		</Select.Group>
	</Select.Content>
</Select.Root>

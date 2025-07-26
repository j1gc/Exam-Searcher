<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';

	let {
		subjectMapping,
		value = $bindable()
	}: { subjectMapping: Map<number, string>; value: number[] } = $props();

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
	<Select.Trigger class="w-[180px] bg-white">
		{#if (subjectIdStrings[0] === '-1' && subjectIdStrings.length === 1) || subjectIdStrings.length === subjectMapping.size}
			Alle Fächer
		{:else if subjectIdStrings.length > 1}
			Mehere Fächer
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

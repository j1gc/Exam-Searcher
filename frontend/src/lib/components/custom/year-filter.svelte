<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { Slider } from '$lib/components/ui/slider/index.js';
	import { m } from '$lib/paraglide/messages';

	let { value = $bindable() }: { value: number[] } = $props();
</script>

<Select.Root type="multiple" name="jahre">
	<Select.Trigger class="w-[100%] bg-white "
		><p class="text-black">{value[0]}-{value[1]}</p></Select.Trigger
	>

	<Select.Content>
		<Select.Group>
			<p class="pt-2 pb-3">{m.year_filter_dropdown()} {value[0]} und {value[1]}</p>
			<Slider
				type="multiple"
				bind:value
				min={2016}
				max={2024}
				step={1}
				class="pb-2"
			/></Select.Group
		>
	</Select.Content>
</Select.Root>

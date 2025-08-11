<script lang="ts">
	import { FileText } from 'lucide-svelte';
	import { Button } from '../ui/button';
	import type { Snippet } from 'svelte';
	import { m } from '$lib/paraglide/messages';

	let { children, fileLink }: { children: Snippet<[]> | undefined; fileLink: string } = $props();
</script>

<div class="grid h-28 grid-cols-4 items-center py-2 max-sm:text-[0.7em] lg:w-[35rem]">
	<div class="col-span-3 flex">
		<div class="flex justify-center px-1 min-sm:px-4">
			<FileText class="text-gray-400 max-sm:hidden min-sm:size-11 " />
		</div>

		<div>
			{@render children?.()}
		</div>
	</div>
	<div>
		<Button href={fileLink} target="_blank" class="max-sm:p-1">{m.file_view_button()}</Button>
	</div>
</div>

<script lang="ts">
	import { FileText } from 'lucide-svelte';
	import { Button } from '../ui/button';
	import type { Snippet } from 'svelte';

	let { children, fileLink }: { children: Snippet<[]> | undefined; fileLink: string } = $props();
</script>

<div class="grid h-28 w-[35rem] grid-cols-4 items-center py-2">
	<div class="col-span-3 flex">
		<div class="flex justify-center px-4">
			<FileText class="size-11 text-gray-400 " />
		</div>

		<div>
			{@render children?.()}
		</div>
	</div>
	<div>
		<Button href={fileLink}>Ansehen</Button>
	</div>
</div>

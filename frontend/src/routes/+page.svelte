<script lang="ts">
	import Button from '$lib/components/ui/button/button.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	let files = data.files;
</script>

<main class="px-52 py-20">
	<h1 class="text-2xl">Files Returned:</h1>
	<div>
		{#each files as file}
			<div class="flex gap-4">
				<p>{file.file_type}</p>
				<p>{file.similarity}</p>
				{#if 'Exam' in file.file}
					<p>{file.file.Exam.file_path}</p>
				{/if}
			</div>
		{/each}
	</div>
</main>

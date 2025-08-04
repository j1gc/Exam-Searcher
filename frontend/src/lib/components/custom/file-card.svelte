<script lang="ts">
	import type { FileResponseSchema } from '$lib/schemas/search';
	import { subjectMapping } from '$lib/subject_mapping.svelte';
	import FileCardOutline from './file-card-outline.svelte';

	let { file }: { file: FileResponseSchema } = $props();

	function filename(file_path: string) {
		return file_path.substring(file_path.lastIndexOf('/') + 1);
	}
</script>

<div class="rounded-md ring-[1.5px] ring-gray-200">
	{#if 'Exam' in file.file}
		<FileCardOutline fileLink={''}>
			<p class="font-semibold">
				{file.file.Exam.difficulty}
				{subjectMapping.get(file.file.Exam.subject_id)} Abitur {file.file.Exam.year}
			</p>
			<p>Prüfung</p>
			<p>{filename(file.file.Exam.file_path)}</p>
		</FileCardOutline>
	{/if}
	{#if 'Answer' in file.file}
		<FileCardOutline fileLink={''}>
			<p class="font-semibold">
				{subjectMapping.get(file.file.Answer.subject_id)}
				{file.file.Answer.year}
			</p>
			<p>Lösung</p>
			<p>{filename(file.file.Answer.file_path)}</p>
		</FileCardOutline>
	{/if}
	{#if 'Other' in file.file}
		<FileCardOutline fileLink={''}>
			<p class="font-semibold">
				{subjectMapping.get(file.file.Other.subject_id)} Abitur {file.file.Other.year}
			</p>
			<p>Anderes</p>
			<p>{filename(file.file.Other.file_path)}</p>
		</FileCardOutline>
	{/if}
</div>

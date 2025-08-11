<script lang="ts">
	import { m } from '$lib/paraglide/messages';
	import type { FileResponseSchema } from '$lib/schemas/search';
	import { subjectMapping } from '$lib/subject_mapping.svelte';
	import FileCardOutline from './file-card-outline.svelte';

	let { file }: { file: FileResponseSchema } = $props();

	function filename(file_path: string) {
		let file = file_path.substring(file_path.lastIndexOf('/') + 1);
		return file.substring(0, file.length - 3) + '.pdf';
	}

	const baseBucketURL = 'https://files.exams.jonas-floerke.de';

	function formatFileURL(filePath: string) {
		// substring removes last 3 chars of path that ends in .md and adds .pdf instead
		let startRemoved = filePath.replace('../../exams/markdown', '');
		let url = baseBucketURL + startRemoved.substring(0, startRemoved.length - 3) + '.pdf';

		return url;
	}
</script>

<div class="rounded-md ring-[1.5px] ring-gray-200">
	{#if 'Exam' in file.file}
		<FileCardOutline fileLink={formatFileURL(file.file.Exam.file_path)}>
			<p class="font-semibold">
				{file.file.Exam.difficulty}
				{subjectMapping.get(file.file.Exam.subject_id)} Abitur {file.file.Exam.year}
			</p>
			<p>{m.file_type_exam()}</p>
			<p>{filename(file.file.Exam.file_path)}</p>
		</FileCardOutline>
	{/if}
	{#if 'Answer' in file.file}
		<FileCardOutline fileLink={formatFileURL(file.file.Answer.file_path)}>
			<p class="font-semibold">
				{subjectMapping.get(file.file.Answer.subject_id)}
				{file.file.Answer.year}
			</p>
			<p>{m.file_type_solution()}</p>
			<p>{filename(file.file.Answer.file_path)}</p>
		</FileCardOutline>
	{/if}
	{#if 'Other' in file.file}
		<FileCardOutline fileLink={file.file.Other.file_path}>
			<p class="font-semibold">
				{subjectMapping.get(file.file.Other.subject_id)} Abitur {file.file.Other.year}
			</p>
			<p>{m.file_type_other()}</p>
			<p>{filename(file.file.Other.file_path)}</p>
		</FileCardOutline>
	{/if}
</div>

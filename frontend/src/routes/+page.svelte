<script lang="ts">
	import * as Form from '$lib/components/ui/form/index.js';

	import { Input } from '$lib/components/ui/input/index.js';
	import { FilesResponseSchema, SearchRequestSchema } from '$lib/schemas/search.js';
	import subjectMappingJson from '../lib/subject_mapping.json';
	import { SubjectMappingSchema } from '$lib/schemas/subject_mapping.js';
	import { get } from 'svelte/store';
	import SubjectFilter from '$lib/components/custom/subject-filter.svelte';
	import FileTypeFilter from '$lib/components/custom/file-type-filter.svelte';

	let subjectMappingData: SubjectMappingSchema = SubjectMappingSchema.parse(subjectMappingJson);
	let subjectMapping = new Map(subjectMappingData.map((subject) => [subject.id, subject.name]));

	async function runQuery(): Promise<FilesResponseSchema> {
		const searchData: SearchRequestSchema = {
			query: searchQuery,
			file_types: selectedFileTypes,
			subject_ids: selectedSubjectIds,
			years: []
		};

		let searchResp = await fetch('http://localhost:3000/search', {
			method: 'POST',
			headers: {
				Accept: 'application/json',
				'Content-Type': 'application/json'
			},

			body: JSON.stringify(searchData)
		});

		return FilesResponseSchema.parse(await searchResp.json());
	}
	let searchQuery = $state('');
	let selectedSubjectIds: number[] = $state([]);
	let selectedFileTypes: string[] = $state(['exam', 'answer', 'other']);

	let returnedFiles: FilesResponseSchema | undefined = $state();

	$effect(() => {
		runQuery().then((r) => (returnedFiles = r));
	});
</script>

<main>
	<h1 class="mb-10 text-3xl font-bold">Finde Prüfungen bei Text, Fach und Jahr</h1>

	<Input placeholder={'Suche bei Text, Fach und Jahr'} bind:value={searchQuery} class="py-7" />

	<div class="flex justify-between pt-5">
		<SubjectFilter {subjectMapping} bind:value={selectedSubjectIds} />
		<FileTypeFilter bind:value={selectedFileTypes} />
		<FileTypeFilter bind:value={selectedFileTypes} />
	</div>
	{#if returnedFiles}
		<h1 class="text-2xl">Files Returned:</h1>
		<div>
			{#each returnedFiles as file}
				<div class="flex gap-4">
					<p>{file.file_type}</p>
					<p>{file.similarity}</p>
					{#if 'Exam' in file.file}
						<p>{file.file.Exam.file_path}</p>
						<p>{subjectMapping.get(file.file.Exam.subject_id)}</p>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</main>

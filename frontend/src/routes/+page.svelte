<script lang="ts">
	import * as Form from '$lib/components/ui/form/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { FilesResponseSchema, SearchRequestSchema } from '$lib/schemas/search.js';

	import SubjectFilter from '$lib/components/custom/subject-filter.svelte';
	import FileTypeFilter from '$lib/components/custom/file-type-filter.svelte';
	import FileCard from '$lib/components/custom/file-card.svelte';
	import YearFilter from '$lib/components/custom/year-filter.svelte';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { subjectMapping } from '$lib/subject_mapping.svelte';
	import Slider from '$lib/components/ui/slider/slider.svelte';

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
	let selectedYears: number[] = $state([2016, 2024]);

	let returnedFiles: FilesResponseSchema | undefined = $state();

	$effect(() => {
		runQuery().then((r) => (returnedFiles = r));
	});
</script>

<main>
	<div class="px-14 py-10">
		<h1 class="mb-10 text-3xl font-bold">Finde Prüfungen und Lösungen bei Text, Fach und Jahr</h1>

		<Input placeholder={'Suche bei Text, Fach und Jahr'} bind:value={searchQuery} class="py-7" />

		<div class="flex justify-between pt-5">
			<SubjectFilter bind:value={selectedSubjectIds} />
			<YearFilter bind:value={selectedYears} />
			<FileTypeFilter bind:value={selectedFileTypes} />
		</div>
	</div>
	<div class="flex bg-white px-14">
		<div class="mr-5 w-72 border-r-1 pt-10 pr-5">
			<p class="text-2xl font-semibold">Angewandte Filter</p>
			<div class="pt-4">
				<p class="font-semibold">Fächer</p>
				<div>
					{#each selectedSubjectIds as id}
						<div class="flex items-center gap-2">
							<Checkbox
								checked={true}
								onclick={(ev) => {
									const idIndex = selectedSubjectIds.indexOf(id);
									selectedSubjectIds.splice(idIndex, 1);
								}}
							/>
							<span>{subjectMapping.get(id)}</span>
						</div>
					{/each}
				</div>
			</div>
			<div class="pt-5">
				<p class="pb-2 font-semibold">Jahre</p>
				<p>{selectedYears[0]}--{selectedYears[1]}</p>
				<Slider type="multiple" bind:value={selectedYears} min={2016} max={2024} step={1} />
			</div>
		</div>

		<div class="pt-10 pl-7">
			{#if returnedFiles}
				<div class="flex flex-col gap-y-5">
					{#each returnedFiles as file}
						<FileCard {file} />
					{/each}
				</div>
			{/if}
		</div>
	</div>
</main>

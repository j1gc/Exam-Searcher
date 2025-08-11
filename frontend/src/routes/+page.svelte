<script lang="ts">
	import { Input } from '$lib/components/ui/input/index.js';
	import { FilesResponseSchema, SearchRequestSchema } from '$lib/schemas/search.js';

	import SubjectFilter from '$lib/components/custom/subject-filter.svelte';
	import FileTypeFilter from '$lib/components/custom/file-type-filter.svelte';
	import FileCard from '$lib/components/custom/file-card.svelte';
	import YearFilter from '$lib/components/custom/year-filter.svelte';
	import Checkbox from '$lib/components/ui/checkbox/checkbox.svelte';
	import { subjectMapping } from '$lib/subject_mapping.svelte';
	import Slider from '$lib/components/ui/slider/slider.svelte';
	import NewestFiles from '$lib/components/custom/newest-files.svelte';
	import { Debounced } from 'runed';
	import { createInfiniteQuery } from '@tanstack/svelte-query';
	import Button from '$lib/components/ui/button/button.svelte';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { m } from '$lib/paraglide/messages';

	async function runQuery(page: number, pageSize: number): Promise<FilesResponseSchema> {
		let years: number[] = [];
		const minYear = Math.min(...selectedYears);
		const maxYear = Math.max(...selectedYears);
		for (let i = minYear; i <= maxYear; i++) {
			years.push(i);
		}

		const searchData: SearchRequestSchema = {
			query: searchQuery,
			file_types: selectedFileTypes,
			subject_ids: selectedSubjectIds,
			years: years
		};

		let searchResp = await fetch(
			`http://localhost:3000/search?index=${page}&page_size=${pageSize}`,
			{
				method: 'POST',
				headers: {
					Accept: 'application/json',
					'Content-Type': 'application/json'
				},

				body: JSON.stringify(searchData)
			}
		);

		return FilesResponseSchema.parse(await searchResp.json());
	}
	let searchQuery = $state('');
	let selectedSubjectIds: number[] = $state([]);
	let selectedFileTypes: string[] = $state(['exam', 'answer', 'other']);
	let selectedYears: number[] = $state([2016, 2024]);

	// only uses the debouncer when one of these values change
	const debouncedHandle = new Debounced(() => ({ searchQuery, selectedYears }), 350);

	const fileQuery = $derived(
		createInfiniteQuery({
			// only executes the query when one of these values change
			queryKey: ['fileQuery', debouncedHandle.current, selectedSubjectIds, selectedFileTypes],
			initialPageParam: { page: 0, pageSize: 20 },

			queryFn: async ({ pageParam }) => {
				return await runQuery(pageParam.page, pageParam.pageSize);
			},
			getNextPageParam: (lastPage, allPages) => {
				// check if last page was a full page
				if (lastPage.length === 20) {
					return { page: allPages.length, pageSize: 20 };
				}
				// Return undefined to indicate no more pages
				return undefined;
			},

			select: (data) => {
				return data.pages.flat();
			}
		})
	);

	let sentinel: HTMLElement;

	onMount(() => {
		const observer = new IntersectionObserver((entries) => {
			if (entries[0].isIntersecting && $fileQuery.hasNextPage) {
				$fileQuery.fetchNextPage();
			}
		});

		observer.observe(sentinel);
	});

	$effect(() => {});
</script>

<main class="pt-5">
	<div class="px-14 pb-10">
		<h1 class="mb-10 text-3xl font-bold text-primary">
			{m.premis_text()}
		</h1>

		<Input placeholder={m.search_field_text()} bind:value={searchQuery} class="py-7" />

		<div class="justify-between gap-x-2 pt-5 min-sm:flex">
			<SubjectFilter bind:value={selectedSubjectIds} />
			<YearFilter bind:value={selectedYears} />
			<FileTypeFilter bind:value={selectedFileTypes} />
		</div>
		<div class="pt-7">
			<h3 class="pb-3 font-semibold">{m.newest_files_text()}</h3>
			<NewestFiles bind:selectedSubject={selectedSubjectIds} />
		</div>
	</div>

	<div class="bg-white px-14 min-sm:flex">
		<div class="mr-5 w-72 border-r-1 pt-10 pr-5">
			<p class="text-2xl font-semibold uppercase">Filter</p>
			<div class="pt-4">
				<p class="font-semibold capitalize">{m.filter_subject_text()}</p>
				<div class="space-y-1">
					{#if selectedSubjectIds.length != 0}{#each selectedSubjectIds as id}
							<div class="flex items-center gap-x-2">
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
					{:else}{#each subjectMapping as subject}
							<div class="flex items-center gap-x-2">
								<Checkbox
									checked={true}
									onclick={(ev) => {
										selectedSubjectIds = Array.from(subjectMapping.keys());

										const idIndex = selectedSubjectIds.indexOf(subject[0]);
										selectedSubjectIds.splice(idIndex, 1);
									}}
								/>
								<span>{subject[1]}</span>
							</div>
						{/each}{/if}
				</div>
			</div>
			<div class="pt-5">
				<p class="pb-2 font-semibold">{m.years_text()}</p>
				<p>{selectedYears[0]}--{selectedYears[1]}</p>
				<Slider type="multiple" bind:value={selectedYears} min={2016} max={2024} step={1} />
			</div>
		</div>

		<div class="pt-10 min-sm:pl-7">
			<h4 class="pb-5 text-2xl font-semibold">{m.results_text()}</h4>
			{#if $fileQuery.isSuccess}
				<div class="flex flex-col gap-y-5">
					{#each $fileQuery.data as file}
						<FileCard {file} />
					{/each}
					<Button
						onclick={$fileQuery.fetchNextPage}
						disabled={!$fileQuery.hasNextPage || $fileQuery.isFetchingNextPage}
					>
						{#if $fileQuery.isFetching}
							Lade...
						{:else if $fileQuery.hasNextPage}
							Lade weitere Prüfungen!
						{:else}
							Keine weiteren Prüfungen gefunden!
						{/if}
					</Button>
				</div>
			{/if}
			<div bind:this={sentinel}></div>
		</div>
	</div>
</main>

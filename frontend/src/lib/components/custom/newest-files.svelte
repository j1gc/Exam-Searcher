<script lang="ts">
	import { BookOpen, BookOpenText, Flower2, Sigma, Omega, Earth } from 'lucide-svelte';
	import { Button } from '../ui/button';
	import * as Carousel from '$lib/components/ui/carousel/index.js';
	import CarouselItem from '../ui/carousel/carousel-item.svelte';
	import { cn } from '$lib/utils';
	import { findSubjectIdByName, subjectMapping } from '$lib/subject_mapping.svelte';

	let { selectedSubject = $bindable() }: { selectedSubject: number[] } = $props();

	let spacing =
		'max-[1499px]:basis-[15%] max-[1270px]:basis-[20%] max-[1000px]:basis-[30%] max-[670px]:basis-[40%] max-[580px]:basis-1/2 max-[500px]:basis-full min-[1500px]:basis-[12%]';

	const newestSubjects = [4, 15, 1, 2, 19, 21];
</script>

<Carousel.Root
	opts={{
		align: 'center'
	}}
>
	<Carousel.Content>
		<CarouselItem class={cn(spacing, '')}>
			<Button
				class="h-40 w-40 rounded-lg bg-[#574D68]"
				onclick={() => (selectedSubject = [newestSubjects[0]])}
			>
				<div class="flex flex-col items-center gap-y-2">
					<Earth class="size-20" />
					<p class="text-lg font-semibold">{subjectMapping.get(newestSubjects[0])}</p>
				</div>
			</Button></CarouselItem
		>
		<Carousel.Item class={cn(spacing, '')}>
			<Button
				class="h-40 w-40 rounded-lg bg-[#8ea4bd]"
				onclick={() => (selectedSubject = [newestSubjects[1]])}
				><div class="flex flex-col items-center gap-y-2">
					<Sigma class="size-20" />
					<p class="text-lg font-semibold">{subjectMapping.get(newestSubjects[1])}</p>
				</div>
			</Button></Carousel.Item
		>
		<Carousel.Item class={cn(spacing, '')}>
			<Button
				class="h-40 w-40 rounded-lg bg-[#e3d3bc] "
				onclick={() => (selectedSubject = [newestSubjects[2]])}
			>
				<div class="flex flex-col items-center gap-y-2">
					<BookOpen class="size-20" />
					<p class="text-lg font-semibold">{subjectMapping.get(newestSubjects[2])}</p>
				</div>
			</Button></Carousel.Item
		>
		<Carousel.Item class={cn(spacing, '')}>
			<Button
				class="h-40 w-40 rounded-lg bg-[#84a297]"
				onclick={() => (selectedSubject = [newestSubjects[3]])}
			>
				<div class="flex flex-col items-center gap-y-2">
					<BookOpenText class="size-20" />
					<p class="text-lg font-semibold">{subjectMapping.get(newestSubjects[3])}</p>
				</div>
			</Button></Carousel.Item
		>

		<CarouselItem class={cn(spacing, '')}>
			<Button
				class="h-40 w-40 rounded-lg bg-[#c8b494]"
				onclick={() => (selectedSubject = [newestSubjects[4]])}
			>
				<div class="flex flex-col items-center gap-y-2">
					<Flower2 class="size-20" />
					<p class="text-lg font-semibold">{subjectMapping.get(newestSubjects[4])}</p>
				</div>
			</Button></CarouselItem
		>
		<CarouselItem class={cn(spacing, '')}>
			<Button
				class="h-40 w-40 rounded-lg bg-[#5F5449]"
				onclick={() => (selectedSubject = [newestSubjects[5]])}
			>
				<div class="flex flex-col items-center gap-y-2">
					<Omega class="size-20" />
					<p class="text-lg font-semibold">{subjectMapping.get(newestSubjects[5])}</p>
				</div>
			</Button></CarouselItem
		>
	</Carousel.Content>
	<Carousel.Previous />
	<Carousel.Next />
</Carousel.Root>

<div class="flex justify-between"></div>

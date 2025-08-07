<script lang="ts">
	import { BookOpen, BookOpenText, Flower2, Sigma, Omega, Earth } from 'lucide-svelte';
	import { Button } from '../ui/button';
	import * as Carousel from '$lib/components/ui/carousel/index.js';
	import CarouselItem from '../ui/carousel/carousel-item.svelte';
	import { cn } from '$lib/utils';

	let spacing =
		'max-[1499px]:basis-[15%] max-[1270px]:basis-[20%] max-[1000px]:basis-[30%] max-[670px]:basis-[40%] max-[580px]:basis-1/2 max-[500px]:basis-full min-[1500px]:basis-[12%]';
</script>

<Carousel.Root
	opts={{
		align: 'start'
	}}
>
	<Carousel.Content>
		<CarouselItem class={cn(spacing, '')}>
			<Button class="h-40 w-40 rounded-lg bg-[#574D68]">
				<div class="flex flex-col items-center gap-y-2">
					<Earth class="size-20" />
					<p class="text-lg font-semibold">Spanisch</p>
				</div>
			</Button></CarouselItem
		>
		<Carousel.Item class={cn(spacing, '')}>
			<Button class="h-40 w-40 rounded-lg bg-[#8ea4bd]"
				><div class="flex flex-col items-center gap-y-2">
					<Sigma class="size-20" />
					<p class="text-lg font-semibold">Mathematik</p>
				</div>
			</Button></Carousel.Item
		>
		<Carousel.Item class={cn(spacing, '')}>
			<Button class="h-40 w-40 rounded-lg bg-[#e3d3bc] ">
				<div class="flex flex-col items-center gap-y-2">
					<BookOpen class="size-20" />
					<p class="text-lg font-semibold">Deutsch</p>
				</div>
			</Button></Carousel.Item
		>
		<Carousel.Item class={cn(spacing, '')}>
			<Button class="h-40 w-40 rounded-lg bg-[#84a297]">
				<div class="flex flex-col items-center gap-y-2">
					<BookOpenText class="size-20" />
					<p class="text-lg font-semibold">Englisch</p>
				</div>
			</Button></Carousel.Item
		>

		<CarouselItem class={cn(spacing, '')}>
			<Button class="h-40 w-40 rounded-lg bg-[#c8b494]">
				<div class="flex flex-col items-center gap-y-2">
					<Flower2 class="size-20" />
					<p class="text-lg font-semibold">Biologie</p>
				</div>
			</Button></CarouselItem
		>
		<CarouselItem class={cn(spacing, '')}>
			<Button class="h-40 w-40 rounded-lg bg-[#5F5449]">
				<div class="flex flex-col items-center gap-y-2">
					<Omega class="size-20" />
					<p class="text-lg font-semibold">Physik</p>
				</div>
			</Button></CarouselItem
		>
	</Carousel.Content>
	<Carousel.Previous />
	<Carousel.Next />
</Carousel.Root>

<div class="flex justify-between"></div>

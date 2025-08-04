<script lang="ts">
	import { BookOpen, BookOpenText, Flower2, Sigma } from 'lucide-svelte';
	import { Button } from '../ui/button';
</script>

<div class="flex justify-between">
	<Button class="h-40 w-40 rounded-lg bg-[#8ea4bd]"
		><div class="flex flex-col items-center gap-y-2">
			<Sigma class="size-20" />
			<p class="text-lg font-semibold">Mathematik</p>
		</div>
	</Button>
	<Button class="h-40 w-40 rounded-lg bg-[#e3d3bc] ">
		<div class="flex flex-col items-center gap-y-2">
			<BookOpen class="size-20" />
			<p class="text-lg font-semibold">Deutsch</p>
		</div>
	</Button>
	<Button class="h-40 w-40 rounded-lg bg-[#84a297]">
		<div class="flex flex-col items-center gap-y-2">
			<BookOpenText class="size-20" />
			<p class="text-lg font-semibold">Englisch</p>
		</div>
	</Button>
	<Button class="h-40 w-40 rounded-lg bg-[#c8b494]">
		<div class="flex flex-col items-center gap-y-2">
			<Flower2 class="size-20" />
			<p class="text-lg font-semibold">Biologie</p>
		</div>
	</Button>
</div>

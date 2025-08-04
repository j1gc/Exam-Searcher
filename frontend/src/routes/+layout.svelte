<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import '../app.css';

	let { children } = $props();
</script>

<div class="h-[100vh] w-[100vw] bg-[#FAF7F3]">
	<div>
		<header class="flex items-center justify-between px-14 py-10">
			<Button href="/" variant="link" class="p-0 text-2xl font-semibold text-black"
				>Prüfungssucher Niedersachsen</Button
			>
			<nav class="flex gap-10">
				<Button href="/" variant="link" class="p-0  text-black">Startseite</Button>
				<Button href="/" variant="link" class="p-0 text-black">Fächer</Button>
				<Button href="/" variant="link" class="p-0 text-black">Jahre</Button>
			</nav>
		</header>
		<main>
			{@render children()}
		</main>
		<footer></footer>
	</div>
</div>

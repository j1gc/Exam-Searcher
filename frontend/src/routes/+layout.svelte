<script lang="ts">
	import '../app.css';

	let { children } = $props();
</script>

<div class="h-[100vh] w-[100vw] bg-[#FAF7F3]">
	<div class="px-20 py-10">
		<header class="mb-20 flex items-center justify-between">
			<h2 class="text-2xl font-semibold">Abituraufgaben Sucher Niedersachsen</h2>
			<nav class="flex gap-10">
				<a href="/">Startseite</a>
				<a href="test">Fächer</a>
				<a href="test">Jahre</a>
			</nav>
		</header>
		<main>
			{@render children()}
		</main>
		<footer></footer>
	</div>
</div>

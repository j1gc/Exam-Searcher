<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import '../app.css';
	import { browser } from '$app/environment';
	import { getLocale, setLocale } from '$lib/paraglide/runtime';
	import { m } from '$lib/paraglide/messages';

	let { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser
			}
		}
	});
</script>

<svelte:head>
	<title>{m.website_title()}</title>
</svelte:head>

<div class="h-[100vh] w-[100vw] bg-[#FAF7F3]">
	<div>
		<header class="items-center justify-between px-14 py-10 min-sm:flex">
			<Button href="/" variant="link" class="p-0 text-2xl font-semibold text-black"
				>{m.website_title()}</Button
			>
			<nav class="flex gap-10">
				<!--<Button href="/" variant="link" class="p-0  text-black">Startseite</Button>
				<Button href="/" variant="link" class="p-0 text-black">Fächer</Button>
				<Button href="/" variant="link" class="p-0 text-black">Jahre</Button>-->
				{#if getLocale() === 'de'}
					<Button
						variant="link"
						class="p-0 text-black"
						onclick={() => {
							setLocale('en');
						}}>Deutsch 🇩🇪</Button
					>
				{:else}<Button
						variant="link"
						class="p-0 text-black"
						onclick={() => {
							setLocale('de');
						}}>English 🇺🇸</Button
					>{/if}
			</nav>
		</header>
		<main>
			<QueryClientProvider client={queryClient}>
				{@render children()}
			</QueryClientProvider>
		</main>
		<footer></footer>
	</div>
</div>

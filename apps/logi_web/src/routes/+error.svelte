<script lang="ts">
	import { page } from '$app/state';
	import { Button } from '$lib/components/ui/button';
	import { onMount } from 'svelte';
</script>

<div class="relative flex h-dvh w-full flex-col items-center justify-center">
	<h1 class="text-7xl font-bold md:text-8xl">{page.status}</h1>
	<p class="mt-0.5 text-2xl text-neutral-600 md:text-3xl dark:text-neutral-400">
		{page.error ? page.error.message : 'An unexpected error occurred.'}
	</p>
	<Button
		class="mt-6"
		onclick={() => {
			// Redirect to the previous page or home
			window.history.back();
		}}
	>
		Go Back
	</Button>
</div>

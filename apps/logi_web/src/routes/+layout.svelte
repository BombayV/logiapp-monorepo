<script lang="ts">
	import '../app.css';
	import { session, updateSession } from '$lib/session';
	import type { LayoutData } from './$types';
	import { Toaster } from '$lib/components/ui/sonner/index.js';

	let { data, children }: { data: LayoutData; children: any } = $props();

	// Update the session store when the layout data changes
	$effect(() => {
		updateSession(data.user);
	});
</script>

<Toaster richColors />

{@render children()}

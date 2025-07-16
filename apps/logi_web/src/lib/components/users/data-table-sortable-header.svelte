<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import ArrowUpDownIcon from '@lucide/svelte/icons/arrow-up-down';
	import ArrowUpIcon from '@lucide/svelte/icons/arrow-up';
	import ArrowDownIcon from '@lucide/svelte/icons/arrow-down';
	import { Button } from '$lib/components/ui/button/index.js';

	let {
		variant = 'ghost',
		sortDirection,
		text,
		...restProps
	}: ComponentProps<typeof Button> & {
		sortDirection?: 'asc' | 'desc' | false;
		text: string;
	} = $props();
</script>

<Button {variant} {...restProps}>
	{text}
	{#if sortDirection === 'asc'}
		<ArrowUpIcon class="ml-2 h-4 w-4" />
	{:else if sortDirection === 'desc'}
		<ArrowDownIcon class="ml-2 h-4 w-4" />
	{:else}
		<ArrowUpDownIcon class="ml-2 h-4 w-4" />
	{/if}
</Button>

<script lang="ts" generics="TData, TValue">
	import {
		getCoreRowModel,
		getSortedRowModel,
		getFilteredRowModel,
		getPaginationRowModel
	} from '@tanstack/table-core';
	import type { ColumnDef, SortingState } from '@tanstack/table-core';
	import { createSvelteTable, FlexRender } from '$lib/components/ui/data-table';
	import * as Table from '$lib/components/ui/table';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Select from '$lib/components/ui/select';
	import { cn } from '$lib/utils';
	import ChevronLeftIcon from '@lucide/svelte/icons/chevron-left';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import ChevronsLeftIcon from '@lucide/svelte/icons/chevrons-left';
	import ChevronsRightIcon from '@lucide/svelte/icons/chevrons-right';

	type DataTableProps<TData, TValue> = {
		columns: ColumnDef<TData, TValue>[];
		data: TData[];
		columnVisibility?: Record<string, boolean>;
		searchable?: boolean;
		searchPlaceholder?: string;
		searchColumn?: string;
		customSearchColumns?: string[];
		paginated?: boolean;
		pageSize?: number;
		sortable?: boolean;
		filterable?: boolean;
		className?: string;
		hideHeader?: boolean;
		hideFooter?: boolean;
		emptyMessage?: string;
	};

	let {
		data,
		columns,
		columnVisibility = {},
		searchable = false,
		searchPlaceholder = 'Search...',
		searchColumn,
		customSearchColumns = [],
		paginated = false,
		pageSize = 10,
		sortable = false,
		filterable = false,
		className,
		hideHeader = false,
		hideFooter = false,
		emptyMessage = 'No results.'
	}: DataTableProps<TData, TValue> = $props();

	let globalFilter = $state('');
	let currentPageSize = $state(pageSize.toString());
	let pagination = $state({
		pageIndex: 0,
		pageSize: pageSize
	});
	let sorting = $state<SortingState>([]);

	// Custom global filter function
	const customGlobalFilter = (row: any, columnId: string, filterValue: string) => {
		if (!filterValue) return true;

		const searchValue = filterValue.toLowerCase();

		// If custom search columns are defined, search only in those columns
		if (customSearchColumns.length > 0) {
			for (const column of customSearchColumns) {
				const cellValue = row.getValue(column);
				if (cellValue && String(cellValue).toLowerCase().includes(searchValue)) {
					return true;
				}
			}
			return false;
		}

		// Default behavior: search all columns
		const allCells = row.getAllCells();
		for (const cell of allCells) {
			const cellValue = cell.getValue();
			if (cellValue && String(cellValue).toLowerCase().includes(searchValue)) {
				return true;
			}
		}
		return false;
	};

	const table = createSvelteTable({
		get data() {
			return data;
		},
		columns,
		getCoreRowModel: getCoreRowModel(),
		getSortedRowModel: sortable ? getSortedRowModel() : undefined,
		getFilteredRowModel: filterable || searchable ? getFilteredRowModel() : undefined,
		getPaginationRowModel: paginated ? getPaginationRowModel() : undefined,
		enableSorting: sortable,
		enableMultiSort: false,
		globalFilterFn: searchable ? customGlobalFilter : undefined,
		onGlobalFilterChange: searchable
			? (updater) => {
					if (typeof updater === 'function') {
						globalFilter = updater(globalFilter);
					} else {
						globalFilter = updater;
					}
				}
			: undefined,
		onPaginationChange: paginated
			? (updater) => {
					if (typeof updater === 'function') {
						pagination = updater(pagination);
					} else {
						pagination = updater;
					}
				}
			: undefined,
		onSortingChange: sortable
			? (updater) => {
					if (typeof updater === 'function') {
						sorting = updater(sorting);
					} else {
						sorting = updater;
					}
				}
			: undefined,
		initialState: {
			columnVisibility: columnVisibility,
			sorting: []
		},
		state: {
			get globalFilter() {
				return searchable ? globalFilter : undefined;
			},
			get pagination() {
				return paginated ? pagination : undefined;
			},
			get sorting() {
				return sortable ? sorting : undefined;
			}
		}
	});
</script>

<div class={cn('w-full', className)}>
	<!-- Search Input -->
	{#if searchable}
		<div class="flex items-center justify-between py-4 gap-x-2">
			<Input
				type="search"
				placeholder={searchPlaceholder}
				bind:value={globalFilter}
				oninput={(e) => {
					table.setGlobalFilter(e.currentTarget.value);
				}}
				class="max-w-sm text-sm"
				autocomplete="off"
				autocorrect="off"
				autocapitalize="off"
				spellcheck="false"
				data-form-type="other"
			/>
			{#if paginated}
				<Select.Root
					type="single"
					bind:value={currentPageSize}
					onValueChange={(value) => {
						if (value && paginated) {
							currentPageSize = value;
							pagination = {
								...pagination,
								pageSize: parseInt(value, 10)
							};
							table.setPageSize(parseInt(value, 10));
						}
					}}
				>
					<Select.Trigger class="h-8 w-[70px]">
						{currentPageSize}
					</Select.Trigger>
					<Select.Content>
						<Select.Label>Numero de filas</Select.Label>
						<Select.Item value="1">1</Select.Item>
						<Select.Item value="10">10</Select.Item>
						<Select.Item value="20">20</Select.Item>
						<Select.Item value="30">30</Select.Item>
						<Select.Item value="40">40</Select.Item>
						<Select.Item value="50">50</Select.Item>
					</Select.Content>
				</Select.Root>
			{/if}
		</div>
	{:else if paginated}
		<div class="flex items-center justify-end py-4">
			<Select.Root
				type="single"
				bind:value={currentPageSize}
				onValueChange={(value) => {
					if (value && paginated) {
						currentPageSize = value;
						pagination = {
							...pagination,
							pageSize: parseInt(value, 10)
						};
						table.setPageSize(parseInt(value, 10));
					}
				}}
			>
				<Select.Trigger class="h-8 w-[70px]">
					{currentPageSize}
				</Select.Trigger>
				<Select.Content>
					<Select.Label>Numero de filas</Select.Label>
					<Select.Item value="1">1</Select.Item>
					<Select.Item value="10">10</Select.Item>
					<Select.Item value="20">20</Select.Item>
					<Select.Item value="30">30</Select.Item>
					<Select.Item value="40">40</Select.Item>
					<Select.Item value="50">50</Select.Item>
				</Select.Content>
			</Select.Root>
		</div>
	{/if}

	<!-- Table -->
	<div class="rounded-md border">
		<Table.Root>
			{#if !hideHeader}
				<Table.Header>
					{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
						<Table.Row>
							{#each headerGroup.headers as header (header.id)}
								<Table.Head colspan={header.colSpan}>
									{#if !header.isPlaceholder}
										<FlexRender
											content={header.column.columnDef.header}
											context={header.getContext()}
										/>
									{/if}
								</Table.Head>
							{/each}
						</Table.Row>
					{/each}
				</Table.Header>
			{/if}
			<Table.Body>
				{@const rowModel = paginated
					? table.getPaginationRowModel()
					: searchable
						? table.getFilteredRowModel()
						: table.getRowModel()}
				{#each rowModel.rows as row (row.id)}
					<Table.Row data-state={row.getIsSelected() && 'selected'}>
						{#each row.getVisibleCells() as cell (cell.id)}
							<Table.Cell>
								<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
							</Table.Cell>
						{/each}
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24 text-center">
							{emptyMessage}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>

	<!-- Pagination -->
	{#if paginated && !hideFooter}
		<div class="flex items-center justify-between space-x-2 py-4">
			<div class="flex-1 text-sm text-muted-foreground">
				{#if table.getFilteredSelectedRowModel().rows.length > 0}
					{table.getFilteredSelectedRowModel().rows.length} de {table.getFilteredRowModel().rows
						.length} filas seleccionadas.
				{:else}
					Mostrando {table.getRowModel().rows.length} de {table.getFilteredRowModel().rows.length} filas.
				{/if}
			</div>
			<div class="flex items-center space-x-6 lg:space-x-8">
				<div class="flex items-center justify-center text-sm font-medium">
					<span class="hidden sm:inline">Pagina&nbsp;</span>
					{table.getState().pagination?.pageIndex + 1}
					<span class="hidden sm:inline whitespace-nowrap">&nbsp;de&nbsp;</span>
					<span class="sm:hidden">/</span>
					{table.getPageCount()}
				</div>
				<div class="flex items-center space-x-2">
					<Button
						variant="outline"
						size="sm"
						onclick={() => table.setPageIndex(0)}
						disabled={!table.getCanPreviousPage()}
						class="hidden sm:inline-flex"
					>
						<span class="sr-only">Ir a la primera página</span>
						<ChevronsLeftIcon class="h-4 w-4" />
					</Button>
					<Button
						variant="outline"
						size="sm"
						onclick={() => table.previousPage()}
						disabled={!table.getCanPreviousPage()}
					>
						<span class="sr-only">Ir a la página anterior</span>
						<ChevronLeftIcon class="h-4 w-4" />
					</Button>
					<Button
						variant="outline"
						size="sm"
						onclick={() => table.nextPage()}
						disabled={!table.getCanNextPage()}
					>
						<span class="sr-only">Ir a la página siguiente</span>
						<ChevronRightIcon class="h-4 w-4" />
					</Button>
					<Button
						variant="outline"
						size="sm"
						onclick={() => table.setPageIndex(table.getPageCount() - 1)}
						disabled={!table.getCanNextPage()}
						class="hidden sm:inline-flex"
					>
						<span class="sr-only">Ir a la última página</span>
						<ChevronsRightIcon class="h-4 w-4" />
					</Button>
				</div>
			</div>
		</div>
	{/if}
</div>

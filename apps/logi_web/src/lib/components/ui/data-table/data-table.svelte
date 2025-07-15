<script lang="ts" generics="TData, TValue">
 import { getCoreRowModel, getSortedRowModel, getFilteredRowModel, getPaginationRowModel } from "@tanstack/table-core";
 import type { ColumnDef } from "@tanstack/table-core";
 import {
  createSvelteTable,
  FlexRender,
 } from "$lib/components/ui/data-table";
 import * as Table from "$lib/components/ui/table";
 import { Button } from "$lib/components/ui/button";
 import { Input } from "$lib/components/ui/input";
 import { cn } from "$lib/utils";
 import ChevronLeftIcon from "@lucide/svelte/icons/chevron-left";
 import ChevronRightIcon from "@lucide/svelte/icons/chevron-right";
 import ChevronsLeftIcon from "@lucide/svelte/icons/chevrons-left";
 import ChevronsRightIcon from "@lucide/svelte/icons/chevrons-right";
 
 type DataTableProps<TData, TValue> = {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  columnVisibility?: Record<string, boolean>;
  searchable?: boolean;
  searchPlaceholder?: string;
  searchColumn?: string;
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
  searchPlaceholder = "Search...",
  searchColumn,
  paginated = false,
  pageSize = 10,
  sortable = false,
  filterable = false,
  className,
  hideHeader = false,
  hideFooter = false,
  emptyMessage = "No results."
 }: DataTableProps<TData, TValue> = $props();
 
 let globalFilter = $state("");
 let currentPageSize = $state(pageSize);
 
 const table = createSvelteTable({
  get data() {
   return data;
  },
  columns,
  getCoreRowModel: getCoreRowModel(),
  getSortedRowModel: sortable ? getSortedRowModel() : undefined,
  getFilteredRowModel: filterable || searchable ? getFilteredRowModel() : undefined,
  getPaginationRowModel: paginated ? getPaginationRowModel() : undefined,
  initialState: {
   columnVisibility: columnVisibility,
   get pagination() {
    return paginated ? {
     pageSize: currentPageSize,
     pageIndex: 0,
    } : undefined;
   },
  },
  get state() {
   return {
    globalFilter: searchable ? globalFilter : undefined,
   };
  },
  onGlobalFilterChange: searchable ? (value) => {
   globalFilter = value;
  } : undefined,
 });
 
 // Handle search column filtering
 $effect(() => {
  if (searchable && searchColumn && globalFilter) {
   table.getColumn(searchColumn)?.setFilterValue(globalFilter);
  } else if (searchable && searchColumn) {
   table.getColumn(searchColumn)?.setFilterValue("");
  }
 });
</script>
 
<div class={cn("w-full", className)}>
 <!-- Search Input -->
 {#if searchable}
  <div class="flex items-center py-4">
   <Input
    placeholder={searchPlaceholder}
    bind:value={globalFilter}
    class="max-w-sm"
   />
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
    {#each table.getRowModel().rows as row (row.id)}
     <Table.Row data-state={row.getIsSelected() && "selected"}>
      {#each row.getVisibleCells() as cell (cell.id)}
       <Table.Cell>
        <FlexRender
         content={cell.column.columnDef.cell}
         context={cell.getContext()}
        />
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
     {table.getFilteredSelectedRowModel().rows.length} of{" "}
     {table.getFilteredRowModel().rows.length} row(s) selected.
    {:else}
     Showing {table.getRowModel().rows.length} of {table.getFilteredRowModel().rows.length} row(s).
    {/if}
   </div>
   <div class="flex items-center space-x-6 lg:space-x-8">
    <div class="flex items-center space-x-2">
     <p class="text-sm font-medium">Rows per page</p>
     <select
      class="flex h-8 w-[70px] items-center justify-center rounded-md border border-input bg-transparent px-2 py-1 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
      bind:value={currentPageSize}
      onchange={() => table.setPageSize(Number(currentPageSize))}
     >
      <option value="10">10</option>
      <option value="20">20</option>
      <option value="30">30</option>
      <option value="40">40</option>
      <option value="50">50</option>
     </select>
    </div>
    <div class="flex w-[100px] items-center justify-center text-sm font-medium">
     Page {table.getState().pagination?.pageIndex + 1} of{" "}
     {table.getPageCount()}
    </div>
    <div class="flex items-center space-x-2">
     <Button
      variant="outline"
      size="sm"
      onclick={() => table.setPageIndex(0)}
      disabled={!table.getCanPreviousPage()}
     >
      <span class="sr-only">Go to first page</span>
      <ChevronsLeftIcon class="h-4 w-4" />
     </Button>
     <Button
      variant="outline"
      size="sm"
      onclick={() => table.previousPage()}
      disabled={!table.getCanPreviousPage()}
     >
      <span class="sr-only">Go to previous page</span>
      <ChevronLeftIcon class="h-4 w-4" />
     </Button>
     <Button
      variant="outline"
      size="sm"
      onclick={() => table.nextPage()}
      disabled={!table.getCanNextPage()}
     >
      <span class="sr-only">Go to next page</span>
      <ChevronRightIcon class="h-4 w-4" />
     </Button>
     <Button
      variant="outline"
      size="sm"
      onclick={() => table.setPageIndex(table.getPageCount() - 1)}
      disabled={!table.getCanNextPage()}
     >
      <span class="sr-only">Go to last page</span>
      <ChevronsRightIcon class="h-4 w-4" />
     </Button>
    </div>
   </div>
  </div>
 {/if}
</div>

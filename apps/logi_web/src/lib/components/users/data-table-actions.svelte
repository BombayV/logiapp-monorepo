<script lang="ts">
 import EllipsisIcon from "@lucide/svelte/icons/ellipsis";
 import { Button } from "$lib/components/ui/button/index.js";
 import * as DropdownMenu from "$lib/components/ui/dropdown-menu/index.js";
 import type { User } from "./columns.js";
 
 let { user }: { user: User } = $props();
 
 const copyEmail = () => {
  navigator.clipboard.writeText(user.email);
  // TODO: Add toast notification
  console.log('Email copiado:', user.email);
 };
 
 const copyPhone = () => {
  if (user.phone) {
   navigator.clipboard.writeText(user.phone);
   // TODO: Add toast notification
   console.log('Teléfono copiado:', user.phone);
  } else {
   alert('No hay teléfono disponible para este usuario');
  }
 };
 
 const deleteUser = () => {
  // TODO: Implement delete user functionality
  console.log('Delete user:', user.user_id);
  alert(`Eliminar usuario: ${user.email}`);
 };
</script>
 
<DropdownMenu.Root>
 <DropdownMenu.Trigger>
  {#snippet child({ props })}
   <Button
    {...props}
    variant="ghost"
    size="icon"
    class="relative size-8 p-0"
   >
    <span class="sr-only">Open menu</span>
    <EllipsisIcon />
   </Button>
  {/snippet}
 </DropdownMenu.Trigger>
 <DropdownMenu.Content>
  <DropdownMenu.Group>
   <DropdownMenu.Label>Acciones</DropdownMenu.Label>
   <DropdownMenu.Item onclick={copyEmail}>
    Copiar email
   </DropdownMenu.Item>
   <DropdownMenu.Item onclick={copyPhone}>
    Copiar teléfono
   </DropdownMenu.Item>
  </DropdownMenu.Group>
  <DropdownMenu.Separator />
  <DropdownMenu.Item variant="destructive" onclick={deleteUser}>
   Borrar usuario
  </DropdownMenu.Item>
 </DropdownMenu.Content>
</DropdownMenu.Root>
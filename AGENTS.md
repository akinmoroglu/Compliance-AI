# AGENTS.md — Rockads Marketing Panel

This file defines the tech stack, conventions, and development guidelines for AI agents and developers working on this project. Follow these rules strictly to ensure consistency with the existing codebase.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Framework | Vue 3 (Composition API) |
| Build Tool | Vite |
| Component Library | shadcn-vue |
| Headless UI Primitives | Reka UI |
| Styling | Tailwind CSS v4 |
| Design Tokens | CSS custom properties via shadcn (oklch color space) |
| Icons | Tabler Icons (`@tabler/icons-vue`) |
| Routing | Vue Router |
| State Management | Pinia (assumed standard for Vue 3 + shadcn-vue projects) |
| Fonts | Google Fonts (via `<link rel="preconnect">`) |
| Analytics | Cloudflare Web Analytics |

---

## Project Conventions

### Vue Components
- Always use the **Composition API** with `<script setup>` syntax.
- Component files use **PascalCase** naming: `UserCard.vue`, `ComplianceTable.vue`.
- Keep components small and single-responsibility. Extract reusable pieces into `/components/ui/` (shadcn primitives) and `/components/` (feature components).
- Use `defineProps` and `defineEmits` with TypeScript types.
```vue
<script setup lang="ts">
const props = defineProps<{
  title: string
  status: 'passed' | 'warning' | 'rejected' | 'pending'
}>()
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle>{{ props.title }}</CardTitle>
    </CardHeader>
  </Card>
</template>
```

---

## shadcn-vue Components

This project uses **shadcn-vue** — the Vue 3 port of shadcn/ui. Components live in `/components/ui/` and are copied into the project (not imported from a package).

### Available Components (already in use)
`Button`, `Card`, `CardHeader`, `CardContent`, `CardDescription`, `Input`, `Badge`, `Table`, `TableHeader`, `TableRow`, `TableHead`, `TableBody`, `TableCell`, `Select`, `SelectTrigger`, `SelectValue`, `Popover`, `PopoverTrigger`, `DropdownMenu`, `DropdownMenuTrigger`, `Sidebar`, `SidebarMenu`, `SidebarMenuButton`, `SidebarMenuItem`, `SidebarGroup`, `SidebarGroupLabel`, `Carousel`, `CarouselContent`, `CarouselItem`, `Avatar`, `AvatarImage`, `Separator`, `Breadcrumb`, `Tooltip`

### Adding New Components
Use the shadcn-vue CLI to add new components:
```bash
npx shadcn-v
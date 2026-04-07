<script setup lang="ts">
import { IconShieldCheck } from '@tabler/icons-vue'
import ComplianceDashboard from './components/ComplianceDashboard.vue'
</script>

<template>
  <div class="min-h-screen bg-background text-foreground flex flex-col">
    <!-- Header -->
    <header class="h-20 border-b border-border bg-card px-8 flex items-center shrink-0">
      <div class="flex items-center gap-3">
        <div class="h-10 w-10 flex items-center justify-center rounded-full bg-primary/10 text-primary">
          <IconShieldCheck :size="24" stroke-width="2" />
        </div>
        <div>
          <h1 class="text-xl font-bold text-primary">Compliance Checker AI</h1>
          <p class="text-sm text-muted-foreground">Verify ad compliance for platform, brand, and legal standards.</p>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 p-8 max-w-5xl mx-auto w-full">
      <ComplianceDashboard />
    </main>
  </div>
</template>

<script setup lang="ts">
import { IconInfoCircle } from '@tabler/icons-vue'

defineProps<{
  advisories: Array<{ code: string; message: string }>
}>()
</script>

<template>
  <div v-if="advisories.length > 0" class="border border-amber-200 rounded-xl bg-amber-50 overflow-hidden">
    <div class="px-4 py-3 border-b border-amber-200 flex items-center gap-2">
      <IconInfoCircle :size="18" class="text-amber-600 shrink-0" />
      <span class="font-bold text-amber-800 text-sm">
        {{ advisories.length }} Recommendation{{ advisories.length > 1 ? 's' : '' }}
      </span>
      <span class="text-xs text-amber-600 ml-1">— these do not affect your compliance score</span>
    </div>
    <div class="p-4 space-y-3">
      <div
        v-for="advisory in advisories"
        :key="advisory.code"
        class="flex items-start gap-2"
      >
        <div class="h-1.5 w-1.5 rounded-full bg-amber-500 mt-2 shrink-0"></div>
        <p class="text-sm text-amber-800">{{ advisory.message }}</p>
      </div>
    </div>
  </div>
</template>

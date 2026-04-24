<script setup lang="ts">
import { ref } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isMobile = ref(window.innerWidth < 640)
window.addEventListener('resize', () => { isMobile.value = window.innerWidth < 640 })

function isActive(path: string) {
  return route.path === path
}
</script>

<template>
  <!-- Desktop: header + main -->
  <div v-if="!isMobile" class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
    <header class="bg-white shadow-sm sticky top-0 z-10">
      <div class="max-w-4xl mx-auto px-4 py-3 flex items-center justify-between">
        <h1 class="text-xl font-bold text-indigo-600">微课脚本生成器</h1>
        <nav class="flex items-center gap-1">
          <button
            v-for="item in [{path:'/',label:'创建'},{path:'/history',label:'历史'},{path:'/settings',label:'设置'}]"
            :key="item.path"
            @click="router.push(item.path)"
            class="px-4 py-2 rounded-lg text-sm font-medium transition"
            :class="isActive(item.path)
              ? 'bg-indigo-100 text-indigo-700'
              : 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'"
          >
            {{ item.label }}
          </button>
        </nav>
      </div>
    </header>
    <main class="max-w-4xl mx-auto px-4 py-8">
      <RouterView />
    </main>
  </div>

  <!-- Mobile: bottom tab bar, no header -->
  <div v-else class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 pb-16">
    <main class="px-4 py-6">
      <RouterView />
    </main>

    <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 z-10">
      <div class="flex">
        <button
          @click="router.push('/')"
          class="flex-1 flex flex-col items-center py-2 gap-0.5 transition"
          :class="isActive('/') ? 'text-indigo-600' : 'text-gray-400'"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
          </svg>
          <span class="text-xs">创建</span>
        </button>

        <button
          @click="router.push('/history')"
          class="flex-1 flex flex-col items-center py-2 gap-0.5 transition"
          :class="isActive('/history') ? 'text-indigo-600' : 'text-gray-400'"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <span class="text-xs">历史</span>
        </button>

        <button
          @click="router.push('/settings')"
          class="flex-1 flex flex-col items-center py-2 gap-0.5 transition"
          :class="isActive('/settings') ? 'text-indigo-600' : 'text-gray-400'"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          <span class="text-xs">设置</span>
        </button>
      </div>
    </nav>
  </div>
</template>

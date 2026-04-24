<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useScriptStore } from '@/stores/script'
import type { Script } from '@/types'

const store = useScriptStore()
const router = useRouter()
const showDeleteModal = ref(false)
const deletingId = ref('')

onMounted(() => {
  store.loadScripts()
})

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return secs > 0 ? `${mins}分${secs}秒` : `${mins}分钟`
}

function styleLabel(style: string): string {
  return style === 'formal' ? '正式' : style === 'casual' ? '轻松' : '讲故事'
}

function styleBadgeColor(style: string): string {
  return style === 'formal' ? 'bg-blue-100 text-blue-700' : style === 'casual' ? 'bg-green-100 text-green-700' : 'bg-orange-100 text-orange-700'
}

function confirmDelete(id: string) {
  deletingId.value = id
  showDeleteModal.value = true
}

async function executeDelete() {
  await store.deleteScript(deletingId.value)
  showDeleteModal.value = false
  deletingId.value = ''
}

function viewScript(script: Script) {
  store.currentScript = script
  router.push('/')
}
</script>

<template>
  <div>
    <!-- 空状态 -->
    <div v-if="store.scripts.length === 0" class="text-center py-20">
      <svg class="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
      </svg>
      <p class="text-gray-400 text-lg mb-4">暂无历史记录</p>
      <button @click="router.push('/')" class="px-6 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition">
        去创建第一个微课
      </button>
    </div>

    <!-- 列表 -->
    <div v-else class="space-y-4">
      <div
        v-for="script in store.scripts"
        :key="script.id"
        class="bg-white rounded-xl shadow-sm hover:shadow-md transition p-5 cursor-pointer"
        @click="viewScript(script)"
      >
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-2 flex-wrap">
              <span class="px-2 py-0.5 rounded text-xs font-medium" :class="styleBadgeColor(script.style)">
                {{ styleLabel(script.style) }}
              </span>
              <span class="text-xs text-gray-400">{{ formatDuration(script.estimatedTotalDuration) }}</span>
              <span v-if="script.targetAudience" class="text-xs text-gray-400">受众: {{ script.targetAudience }}</span>
            </div>
            <h3 class="text-lg font-semibold text-gray-800 truncate">{{ script.topic }}</h3>
            <p class="text-sm text-gray-400 mt-1">{{ formatDate(script.createdAt) }}</p>
            <p class="text-xs text-gray-400 mt-1">{{ script.sections.length }} 个章节</p>
          </div>
          <div class="flex flex-col gap-1" @click.stop>
            <button
              @click="viewScript(script)"
              class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition"
              title="查看"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
              </svg>
            </button>
            <button
              @click="confirmDelete(script.id)"
              class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition"
              title="删除"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 px-4">
        <div class="bg-white rounded-xl shadow-xl p-6 max-w-sm w-full">
          <h3 class="text-lg font-semibold text-gray-800 mb-2">确认删除</h3>
          <p class="text-gray-500 text-sm mb-6">删除后将无法恢复，确定要删除吗？</p>
          <div class="flex gap-3 justify-end">
            <button @click="showDeleteModal = false" class="px-4 py-2 text-gray-600 hover:text-gray-800 transition text-sm">
              取消
            </button>
            <button @click="executeDelete" class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition text-sm">
              删除
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Script } from '@/types'

const props = defineProps<{
  script: Script | null
  previewText?: string
}>()

const editingSection = ref<string | null>(null)
const editContent = ref('')

function startEdit(sectionId: string, content: string) {
  editingSection.value = sectionId
  editContent.value = content
}

function saveEdit() {
  if (editingSection.value) {
    const section = props.script.sections.find(s => s.id === editingSection.value)
    if (section) {
      section.content = editContent.value
    }
    editingSection.value = null
  }
}

function cancelEdit() {
  editingSection.value = null
  editContent.value = ''
}

function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}秒`
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return secs > 0 ? `${mins}分${secs}秒` : `${mins}分钟`
}

function copyToClipboard(text: string) {
  navigator.clipboard.writeText(text)
}

function exportMarkdown(): string {
  if (!props.script) return ''
  let md = `# ${props.script.topic}\n\n`
  md += `**目标受众**: ${props.script.targetAudience}  |  **风格**: ${props.script.style}\n\n`
  md += `---\n\n`

  for (const section of props.script.sections) {
    md += `## ${section.title}\n\n`
    md += `> 时长: ${formatDuration(section.estimatedDuration)}\n\n`
    md += `${section.content}\n\n`
  }

  return md
}

function downloadMarkdown() {
  if (!props.script) return
  const md = exportMarkdown()
  const blob = new Blob([md], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${props.script.topic}.md`
  a.click()
  URL.revokeObjectURL(url)
}

function sectionTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    opening: '开场',
    knowledge: '知识点',
    example: '案例',
    summary: '总结'
  }
  return labels[type] || type
}

function sectionTypeColor(type: string): string {
  const colors: Record<string, string> = {
    opening: 'bg-green-100 text-green-700',
    knowledge: 'bg-blue-100 text-blue-700',
    example: 'bg-yellow-100 text-yellow-700',
    summary: 'bg-purple-100 text-purple-700'
  }
  return colors[type] || 'bg-gray-100 text-gray-700'
}
</script>

<template>
  <div class="bg-white rounded-xl shadow-lg p-6">
    <!-- 正在生成中 -->
    <div v-if="!script && previewText" class="space-y-4">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center">
          <svg class="animate-spin w-5 h-5 text-indigo-600" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
          </svg>
        </div>
        <div>
          <p class="font-semibold text-gray-800">正在生成微课脚本...</p>
          <p class="text-xs text-gray-400">AI 正在构思中，请稍候</p>
        </div>
      </div>
      <div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
        <pre class="text-sm text-gray-700 whitespace-pre-wrap font-sans leading-relaxed">{{ previewText }}</pre>
      </div>
    </div>

    <!-- 脚本已生成 -->
    <template v-else-if="script">
    <!-- Header -->
    <div class="flex items-start justify-between mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-800 mb-2">{{ script.topic }}</h2>
        <div class="flex items-center gap-4 text-sm text-gray-500">
          <span>预计时长: {{ formatDuration(script.estimatedTotalDuration) }}</span>
          <span>风格: {{ script.style === 'formal' ? '正式' : script.style === 'casual' ? '轻松' : '讲故事' }}</span>
        </div>
      </div>
      <div class="flex gap-2">
        <button
          @click="downloadMarkdown"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition text-sm"
        >
          导出 Markdown
        </button>
      </div>
    </div>

    <!-- Sections -->
    <div class="space-y-4">
      <div
        v-for="section in script.sections"
        :key="section.id"
        class="border border-gray-200 rounded-lg p-4"
      >
        <div class="flex items-center justify-between mb-3">
          <div class="flex items-center gap-3">
            <span
              class="px-2 py-1 rounded text-xs font-medium"
              :class="sectionTypeColor(section.type)"
            >
              {{ sectionTypeLabel(section.type) }}
            </span>
            <h3 class="font-semibold text-gray-800">{{ section.title }}</h3>
            <span class="text-xs text-gray-400">
              {{ formatDuration(section.estimatedDuration) }}
            </span>
          </div>
          <div class="flex gap-1">
            <button
              v-if="editingSection !== section.id"
              @click="startEdit(section.id, section.content)"
              class="p-1.5 text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 rounded transition"
              title="编辑"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
              </svg>
            </button>
            <button
              @click="copyToClipboard(section.content)"
              class="p-1.5 text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 rounded transition"
              title="复制"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- Content Display -->
        <div v-if="editingSection !== section.id">
          <p class="text-gray-700 whitespace-pre-wrap leading-relaxed">
            {{ section.content }}
          </p>
        </div>

        <!-- Edit Mode -->
        <div v-else class="space-y-3">
          <textarea
            v-model="editContent"
            rows="8"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none resize-none font-mono text-sm"
          />
          <div class="flex gap-2 justify-end">
            <button
              @click="cancelEdit"
              class="px-4 py-1.5 text-gray-600 hover:text-gray-800 transition text-sm"
            >
              取消
            </button>
            <button
              @click="saveEdit"
              class="px-4 py-1.5 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition text-sm"
            >
              保存
            </button>
          </div>
        </div>
      </div>
    </div>
    </template>
  </div>
</template>

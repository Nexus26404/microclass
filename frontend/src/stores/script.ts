import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Script, GenerateRequest } from '@/types'

export const useScriptStore = defineStore('script', () => {
  const scripts = ref<Script[]>([])
  const currentScript = ref<Script | null>(null)
  const previewText = ref('')
  const isGenerating = ref(false)
  const error = ref<string | null>(null)

  function setScript(script: Script) {
    currentScript.value = script
  }

  function setPreview(text: string) {
    previewText.value = text
  }

  async function generateScript(request: GenerateRequest): Promise<void> {
    if (isGenerating.value) return

    isGenerating.value = true
    error.value = null
    previewText.value = ''

    try {
      const response = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request)
      })

      if (!response.ok) {
        let msg = '生成失败'
        try {
          const data = await response.json()
          msg = data.error || msg
        } catch {}
        throw new Error(msg)
      }

      const data = await response.json()
      currentScript.value = data.script
    } catch (e) {
      error.value = e instanceof Error ? e.message : '生成失败'
    } finally {
      isGenerating.value = false
    }
  }

  async function loadScripts() {
    try {
      const response = await fetch('/api/scripts')
      const data = await response.json()
      scripts.value = data.scripts || []
    } catch (e) {
      console.error('加载历史失败:', e)
    }
  }

  async function deleteScript(id: string) {
    await fetch(`/api/scripts/${id}`, { method: 'DELETE' })
    scripts.value = scripts.value.filter(s => s.id !== id)
    if (currentScript.value?.id === id) {
      currentScript.value = null
    }
  }

  function updateSection(sectionId: string, content: string) {
    if (currentScript.value) {
      const section = currentScript.value.sections.find(s => s.id === sectionId)
      if (section) {
        section.content = content
      }
    }
  }

  return {
    scripts,
    currentScript,
    previewText,
    isGenerating,
    error,
    setScript,
    setPreview,
    generateScript,
    loadScripts,
    deleteScript,
    updateSection
  }
})

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Settings {
  mockMode: boolean
  openAIURL: string
  openAIModel: string
}

const settings = ref<Settings>({
  mockMode: true,
  openAIURL: '',
  openAIModel: ''
})

const apiKey = ref('')
const isLoading = ref(true)
const isSaving = ref(false)
const saveMsg = ref('')
const saveMsgType = ref<'success' | 'error'>('success')

onMounted(async () => {
  try {
    const res = await fetch('/api/settings')
    const data = await res.json()
    settings.value = {
      mockMode: data.settings?.mockMode ?? true,
      openAIURL: data.settings?.openAIURL ?? '',
      openAIModel: data.settings?.openAIModel ?? ''
    }
  } catch {
    showMsg('加载设置失败', 'error')
  } finally {
    isLoading.value = false
  }
})

function showMsg(msg: string, type: 'success' | 'error') {
  saveMsg.value = msg
  saveMsgType.value = type
  setTimeout(() => { saveMsg.value = '' }, 3000)
}

async function saveSettings() {
  isSaving.value = true
  try {
    const body: Record<string, unknown> = {
      mockMode: settings.value.mockMode
    }
    if (settings.value.openAIURL.trim()) {
      body.openAIURL = settings.value.openAIURL.trim()
    }
    if (settings.value.openAIModel.trim()) {
      body.openAIModel = settings.value.openAIModel.trim()
    }
    if (apiKey.value.trim()) {
      body.openAIKey = apiKey.value.trim()
      apiKey.value = ''
    }

    const res = await fetch('/api/settings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body)
    })
    if (!res.ok) throw new Error('保存失败')
    showMsg('设置已保存', 'success')
  } catch (e) {
    showMsg(e instanceof Error ? e.message : '保存失败', 'error')
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div>
    <div v-if="isLoading" class="text-center py-12 text-gray-400">加载中...</div>

    <div v-else class="space-y-6">
        <!-- 提示卡片 -->
        <div class="bg-amber-50 border border-amber-200 rounded-xl p-4 text-sm text-amber-700">
          <p>修改设置后会立即生效。API Key 仅在保存时传递，不会明文显示。</p>
        </div>

        <!-- 生成模式 -->
        <div class="bg-white rounded-xl shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-1">生成模式</h2>
          <p class="text-sm text-gray-500 mb-4">切换生成脚本的方式</p>

          <div class="space-y-3">
            <label class="flex items-center gap-3 cursor-pointer group">
              <input type="radio" :value="true" v-model="settings.mockMode" class="accent-indigo-600 w-4 h-4" />
              <div>
                <p class="font-medium text-gray-800 group-hover:text-indigo-600 transition">模拟模式</p>
                <p class="text-xs text-gray-400">使用本地模拟数据，无需 API Key，适合测试界面</p>
              </div>
            </label>
            <label class="flex items-center gap-3 cursor-pointer group">
              <input type="radio" :value="false" v-model="settings.mockMode" class="accent-indigo-600 w-4 h-4" />
              <div>
                <p class="font-medium text-gray-800 group-hover:text-indigo-600 transition">真实 API</p>
                <p class="text-xs text-gray-400">调用 MiniMax GPT 等 OpenAI 兼容接口，生成真实内容</p>
              </div>
            </label>
          </div>
        </div>

        <!-- API 配置 -->
        <div class="bg-white rounded-xl shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-1">API 配置</h2>
          <p class="text-sm text-gray-500 mb-4">配置 OpenAI 兼容接口参数</p>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">API URL</label>
              <input
                v-model="settings.openAIURL"
                type="url"
                placeholder="https://api.minimax.chat"
                class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none text-sm transition"
              />
              <p class="text-xs text-gray-400 mt-1">OpenAI 兼容格式的 API 地址</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">模型名称</label>
              <input
                v-model="settings.openAIModel"
                type="text"
                placeholder="MiniMax-Text-01"
                class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none text-sm transition"
              />
              <p class="text-xs text-gray-400 mt-1">例如：MiniMax-Text-01、gpt-4o</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                API Key <span class="text-xs text-gray-400 font-normal">（留空则不修改）</span>
              </label>
              <input
                v-model="apiKey"
                type="password"
                placeholder="输入新密钥以更新"
                class="w-full px-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent outline-none text-sm transition"
              />
              <p class="text-xs text-gray-400 mt-1">API Key 不会明文回显，仅在保存时提交</p>
            </div>
          </div>
        </div>

        <!-- 保存按钮 -->
        <button
          @click="saveSettings"
          :disabled="isSaving"
          class="w-full py-3 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 active:bg-indigo-800 transition disabled:opacity-50"
        >
          <span v-if="isSaving" class="flex items-center justify-center gap-2">
            <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
            保存中...
          </span>
          <span v-else>保存设置</span>
        </button>

        <!-- 消息提示 -->
        <div v-if="saveMsg" class="text-center text-sm" :class="saveMsgType === 'success' ? 'text-green-600' : 'text-red-500'">
          {{ saveMsg }}
        </div>
      </div>
  </div>
</template>

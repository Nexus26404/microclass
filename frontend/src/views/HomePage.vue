<script setup lang="ts">
import { ref, computed, reactive } from 'vue'
import { useScriptStore } from '@/stores/script'
import ScriptDisplay from '@/components/ScriptDisplay.vue'

const store = useScriptStore()

const defaultForm = () => ({
  topic: '',
  style: 'formal' as const,
  showAdvanced: false,
  discipline: '',
  targetAudience: '',
  targetDuration: 180,
  customStructure: '',
  referenceMode: 'paste' as const,
  referenceText: '',
})

const topic = ref('')
const style = ref<'formal' | 'casual' | 'storytelling'>('formal')
const showAdvanced = ref(false)
const discipline = ref('')
const targetAudience = ref('')
const targetDuration = ref(180)
const customStructure = ref('')

const referenceMode = ref<'paste' | 'upload'>('paste')
const referenceText = ref('')
const referenceFile = ref<File | null>(null)
const referenceFileName = ref('')
const isUploading = ref(false)
const uploadError = ref('')

const formErrors = reactive({
  topic: '',
  discipline: '',
  targetAudience: '',
  customStructure: '',
  referenceText: ''
})

function validateTopic() {
  const v = topic.value.trim()
  if (!v) {
    formErrors.topic = '请输入微课主题'
  } else if (v.length < 2) {
    formErrors.topic = '主题至少需要 2 个字符'
  } else if (v.length > 200) {
    formErrors.topic = '主题不能超过 200 个字符'
  } else {
    formErrors.topic = ''
  }
}

function validateDiscipline() {
  formErrors.discipline = discipline.value.trim().length > 100 ? '学科/行业不能超过 100 个字符' : ''
}

function validateAudience() {
  formErrors.targetAudience = targetAudience.value.trim().length > 100 ? '目标受众不能超过 100 个字符' : ''
}

function validateStructure() {
  formErrors.customStructure = customStructure.value.trim().length > 500 ? '脚本结构不能超过 500 个字符' : ''
}

function validateReferenceText() {
  if (referenceText.value.length > 10000) {
    formErrors.referenceText = '粘贴文本不能超过 10000 个字符'
  } else {
    formErrors.referenceText = ''
  }
}

const canGenerate = computed(() => {
  validateTopic()
  return topic.value.trim().length >= 2 && !store.isGenerating && !formErrors.topic && !formErrors.referenceText
})

function handleFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  uploadError.value = ''
  const ext = file.name.split('.').pop()?.toLowerCase()
  if (ext !== 'txt' && ext !== 'md') {
    uploadError.value = '仅支持 .txt 和 .md 格式文件'
    referenceFile.value = null
    referenceFileName.value = ''
    return
  }
  if (file.size > 2 * 1024 * 1024) {
    uploadError.value = '文件大小不能超过 2MB'
    referenceFile.value = null
    referenceFileName.value = ''
    return
  }

  referenceFile.value = file
  referenceFileName.value = file.name

  isUploading.value = true
  const formData = new FormData()
  formData.append('file', file)

  fetch('/api/upload', { method: 'POST', body: formData })
    .then(res => res.json())
    .then(data => {
      if (data.error) {
        uploadError.value = data.error
        referenceFile.value = null
        referenceFileName.value = ''
      } else {
        referenceText.value = data.content
      }
    })
    .catch(() => {
      uploadError.value = '上传失败，请重试'
      referenceFile.value = null
      referenceFileName.value = ''
    })
    .finally(() => {
      isUploading.value = false
    })
}

function removeFile() {
  referenceFile.value = null
  referenceFileName.value = ''
  referenceText.value = ''
  uploadError.value = ''
}

const abortController = ref<AbortController | null>(null)

async function handleGenerate() {
  validateTopic()
  validateDiscipline()
  validateAudience()
  validateStructure()
  validateReferenceText()

  if (formErrors.topic || formErrors.discipline || formErrors.targetAudience || formErrors.customStructure || formErrors.referenceText) {
    return
  }

  abortController.value = new AbortController()
  store.isGenerating = true
  store.error = null
  store.previewText = ''

  try {
    const response = await fetch('/api/generate-stream', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        topic: topic.value.trim(),
        style: style.value,
        discipline: discipline.value.trim() || undefined,
        targetAudience: targetAudience.value.trim() || undefined,
        targetDuration: targetDuration.value || undefined,
        customStructure: customStructure.value.trim() || undefined,
        referenceText: referenceText.value.trim() || undefined
      }),
      signal: abortController.value.signal
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || '生成失败')
    }

    const reader = response.body!.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() ?? ''

      for (const line of lines) {
        if (!line.startsWith('event:') && !line.startsWith('data:')) continue
        if (line.startsWith('event:')) continue
        const data = line.slice(5).trim()
        if (!data) continue

        try {
          const parsed = JSON.parse(data)
          if (parsed.text !== undefined) {
            store.previewText += parsed.text
          }
          if (parsed.script) {
            store.setScript(parsed.script)
          }
          if (parsed.error) {
            throw new Error(parsed.error)
          }
        } catch {}
      }
    }
  } catch (e) {
    if (e instanceof DOMException && e.name === 'AbortError') {
      store.error = null
    } else {
      store.error = e instanceof Error ? e.message : '生成失败'
    }
  } finally {
    store.isGenerating = false
    abortController.value = null
  }
}

function cancelGenerate() {
  abortController.value?.abort()
}

const topicInputClasses = computed(() =>
  formErrors.topic ? 'border-red-400 focus:ring-red-500' : 'border-gray-300 focus:ring-indigo-500'
)
const disciplineInputClasses = computed(() =>
  formErrors.discipline ? 'border-red-400 focus:ring-red-500' : 'border-gray-300 focus:ring-indigo-500'
)
const audienceInputClasses = computed(() =>
  formErrors.targetAudience ? 'border-red-400 focus:ring-red-500' : 'border-gray-300 focus:ring-indigo-500'
)
const structureInputClasses = computed(() =>
  formErrors.customStructure ? 'border-red-400 focus:ring-red-500' : 'border-gray-300 focus:ring-indigo-500'
)
const refTextInputClasses = computed(() =>
  formErrors.referenceText ? 'border-red-400 focus:ring-red-500' : 'border-gray-300 focus:ring-indigo-500'
)

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return secs === 0 ? `${mins}分钟` : `${mins}分${secs}秒`
}

function resetForm() {
  const d = defaultForm()
  topic.value = d.topic
  style.value = d.style
  showAdvanced.value = d.showAdvanced
  discipline.value = d.discipline
  targetAudience.value = d.targetAudience
  targetDuration.value = d.targetDuration
  customStructure.value = d.customStructure
  referenceMode.value = d.referenceMode
  referenceText.value = d.referenceText
  referenceFile.value = null
  referenceFileName.value = ''
  uploadError.value = ''
  store.error = null
  store.currentScript = null
  Object.keys(formErrors).forEach(k => (formErrors[k as keyof typeof formErrors] = ''))
}
</script>

<template>
  <div>
    <div class="bg-white rounded-xl shadow-lg p-6 mb-6">
      <h2 class="text-xl font-semibold text-gray-800 mb-4">创建新微课</h2>

        <div class="space-y-4">
          <!-- 微课主题 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">
              微课主题 <span class="text-red-500">*</span>
            </label>
            <input
              v-model="topic"
              type="text"
              placeholder="例如：Python 列表推导式的使用技巧"
              class="w-full px-4 py-3 border rounded-lg focus:ring-2 focus:border-transparent outline-none transition"
              :class="topicInputClasses"
              @blur="validateTopic"
              @input="formErrors.topic = ''"
              @keyup.enter="handleGenerate"
            />
            <p v-if="formErrors.topic" class="mt-1 text-sm text-red-500">{{ formErrors.topic }}</p>
          </div>

          <!-- 风格 -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              风格 <span class="text-red-500">*</span>
            </label>
            <div class="flex gap-3">
              <label
                v-for="s in [
                  { value: 'formal', label: '正式', desc: '严谨专业' },
                  { value: 'casual', label: '轻松', desc: '口语化' },
                  { value: 'storytelling', label: '讲故事', desc: '引人入胜' }
                ]"
                :key="s.value"
                class="flex-1 cursor-pointer"
              >
                <input type="radio" :value="s.value" v-model="style" class="sr-only peer" />
                <div class="px-4 py-3 border rounded-lg text-center transition peer-checked:border-indigo-500 peer-checked:bg-indigo-50 peer-checked:text-indigo-700 hover:border-indigo-300">
                  <div class="font-medium">{{ s.label }}</div>
                  <div class="text-xs text-gray-500">{{ s.desc }}</div>
                </div>
              </label>
            </div>
          </div>

          <!-- 高级选项折叠区 -->
          <div class="border border-gray-200 rounded-lg overflow-hidden">
            <button
              @click="showAdvanced = !showAdvanced"
              class="w-full px-4 py-3 flex items-center justify-between text-sm text-gray-600 hover:bg-gray-50 transition"
            >
              <span>高级选项</span>
              <svg class="w-4 h-4 transition-transform" :class="showAdvanced ? 'rotate-180' : ''" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
              </svg>
            </button>

            <div v-show="showAdvanced" class="px-4 pb-4 space-y-4 border-t border-gray-100 pt-4">
              <!-- 学科/行业 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">学科 / 行业</label>
                <input
                  v-model="discipline"
                  type="text"
                  placeholder="例如：编程、设计、产品运营"
                  class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:border-transparent outline-none text-sm transition"
                  :class="disciplineInputClasses"
                  @blur="validateDiscipline"
                  @input="formErrors.discipline = ''"
                />
                <p v-if="formErrors.discipline" class="mt-1 text-sm text-red-500">{{ formErrors.discipline }}</p>
              </div>

              <!-- 目标受众 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">目标受众</label>
                <input
                  v-model="targetAudience"
                  type="text"
                  placeholder="例如：零基础学员、大学生、职场新人"
                  class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:border-transparent outline-none text-sm transition"
                  :class="audienceInputClasses"
                  @blur="validateAudience"
                  @input="formErrors.targetAudience = ''"
                />
                <p v-if="formErrors.targetAudience" class="mt-1 text-sm text-red-500">{{ formErrors.targetAudience }}</p>
              </div>

              <!-- 预计时长 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  预计时长 <span class="text-xs text-gray-400 font-normal">（滑动选择，最长 15 分钟）</span>
                </label>
                <div class="flex items-center gap-3">
                  <input
                    v-model.number="targetDuration"
                    type="range" min="60" max="900" step="30" class="flex-1"
                  />
                  <span class="text-sm font-medium text-indigo-600 w-20 text-right">
                    {{ formatDuration(targetDuration) }}
                  </span>
                </div>
                <div class="flex justify-between text-xs text-gray-400 mt-1">
                  <span>1分钟</span><span>15分钟</span>
                </div>
              </div>

              <!-- 脚本结构 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">脚本结构 <span class="text-gray-400 font-normal">（可选）</span></label>
                <input
                  v-model="customStructure"
                  type="text"
                  placeholder="例如：[话题导入]----[知识点讲解]----[习题]----[总结]"
                  class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:border-transparent outline-none text-sm transition"
                  :class="structureInputClasses"
                  @blur="validateStructure"
                  @input="formErrors.customStructure = ''"
                />
                <p v-if="formErrors.customStructure" class="mt-1 text-sm text-red-500">{{ formErrors.customStructure }}</p>
                <p class="text-xs text-gray-400 mt-1">用 ---- 分隔各段落，如留空则使用默认结构</p>
              </div>

              <!-- 参考资料 -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">参考资料 <span class="text-gray-400 font-normal">（可选，有助于生成更精准的脚本）</span></label>

                <!-- 切换标签 -->
                <div class="flex gap-1 mb-2">
                  <button
                    @click="referenceMode = 'paste'"
                    class="px-3 py-1.5 text-sm rounded-md transition"
                    :class="referenceMode === 'paste'
                      ? 'bg-indigo-100 text-indigo-700 font-medium'
                      : 'text-gray-500 hover:bg-gray-100'"
                  >
                    粘贴文本
                  </button>
                  <button
                    @click="referenceMode = 'upload'"
                    class="px-3 py-1.5 text-sm rounded-md transition"
                    :class="referenceMode === 'upload'
                      ? 'bg-indigo-100 text-indigo-700 font-medium'
                      : 'text-gray-500 hover:bg-gray-100'"
                  >
                    上传文件
                  </button>
                </div>

                <!-- 粘贴模式 -->
                <div v-if="referenceMode === 'paste'">
                  <textarea
                    v-model="referenceText"
                    rows="5"
                    placeholder="粘贴参考资料内容，例如文章、笔记、教材摘要等..."
                    class="w-full px-3 py-2 border rounded-lg focus:ring-2 focus:border-transparent outline-none text-sm resize-none transition"
                    :class="refTextInputClasses"
                    @blur="validateReferenceText"
                    @input="formErrors.referenceText = ''"
                  ></textarea>
                  <p v-if="formErrors.referenceText" class="mt-1 text-sm text-red-500">{{ formErrors.referenceText }}</p>
                  <p class="text-xs text-gray-400 mt-1">建议控制字数，过长的内容会自动截断</p>
                </div>

                <!-- 上传模式 -->
                <div v-if="referenceMode === 'upload'">
                  <div v-if="!referenceFileName">
                    <label class="block w-full px-4 py-6 border-2 border-dashed border-gray-300 rounded-lg text-center cursor-pointer hover:border-indigo-400 hover:bg-indigo-50 transition">
                      <span class="text-indigo-600 font-medium">点击上传文件</span>
                      <p class="text-xs text-gray-400 mt-1">支持 .txt、.md 格式，最大 2MB</p>
                      <input
                        type="file"
                        accept=".txt,.md"
                        class="hidden"
                        @change="handleFileChange"
                      />
                    </label>
                  </div>

                  <!-- 已上传文件 -->
                  <div v-else class="flex items-center justify-between px-4 py-3 bg-green-50 border border-green-200 rounded-lg">
                    <div class="flex items-center gap-2">
                      <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                      </svg>
                      <div>
                        <p class="text-sm font-medium text-green-800">{{ referenceFileName }}</p>
                        <p class="text-xs text-green-600">{{ referenceText.length }} 个字符</p>
                      </div>
                    </div>
                    <button @click="removeFile" class="text-green-600 hover:text-green-800 p-1">
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                      </svg>
                    </button>
                  </div>

                  <p v-if="uploadError" class="mt-1 text-sm text-red-500">{{ uploadError }}</p>
                  <p v-if="isUploading" class="mt-1 text-sm text-indigo-500">正在解析文件...</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex gap-3">
            <button
              v-if="store.isGenerating"
              @click="cancelGenerate"
              class="flex-1 py-3 bg-red-500 text-white rounded-lg font-medium hover:bg-red-600 active:bg-red-700 transition"
            >
              <span class="flex items-center justify-center gap-2">
                <svg class="animate-spin h-5 w-5" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"/>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
                </svg>
                取消生成
              </span>
            </button>
            <button
              v-else
              @click="handleGenerate"
              :disabled="!canGenerate"
              class="flex-1 py-3 bg-indigo-600 text-white rounded-lg font-medium transition"
              :class="canGenerate ? 'hover:bg-indigo-700 active:bg-indigo-800' : 'opacity-50 cursor-not-allowed'"
            >
              生成微课脚本
            </button>
            <button
              @click="resetForm"
              class="px-6 py-3 border border-gray-300 text-gray-600 rounded-lg font-medium hover:bg-gray-50 transition"
              :class="store.isGenerating ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-50'"
              :disabled="store.isGenerating"
            >
              重置
            </button>
          </div>
        </div>
      </div>

      <!-- 错误提示 -->
      <div v-if="store.error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
        {{ store.error }}
      </div>

      <!-- 脚本展示 -->
      <ScriptDisplay v-if="store.currentScript || store.previewText" :script="store.currentScript" :preview-text="store.previewText" />
  </div>
</template>

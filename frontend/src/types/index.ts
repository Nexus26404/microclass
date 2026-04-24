export interface ScriptSection {
  id: string
  type: string
  title: string
  content: string
  estimatedDuration: number
}

export interface Script {
  id: string
  topic: string
  targetAudience: string
  style: 'formal' | 'casual' | 'storytelling'
  estimatedTotalDuration: number
  sections: ScriptSection[]
  createdAt: string
  updatedAt: string
}

export interface GenerateRequest {
  topic: string
  discipline?: string
  targetAudience?: string
  targetDuration?: number
  style: 'formal' | 'casual' | 'storytelling'
  customStructure?: string
  referenceText?: string
}

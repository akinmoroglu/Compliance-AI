<script setup lang="ts">
import { ref, computed } from 'vue'
import { 
  IconCheck, 
  IconAlertCircle, 
  IconUpload, 
  IconChevronDown, 
  IconBrandMeta, 
  IconBrandTiktok, 
  IconBrandGoogle,
  IconLoader2,
  IconUsers
} from '@tabler/icons-vue'
import AdvisoryPanel from './AdvisoryPanel.vue'

const currentStep = ref(1)
const isChecking = ref(false)
const isSubmitting = ref(false)

// Progress UI state simulation
const timeElapsed = ref(0)
let timerInterval: any = null
const loadingSteps = ref([
  { label: 'Initializing compliance scan...', status: 'done' },
  { label: 'Uploading & Preprocessing media...', status: 'pending' },
  { label: 'Classifying ad category...', status: 'pending' },
  { label: 'Evaluating wording and visual creatives...', status: 'pending' },
  { label: 'Analyzing web claims & alignments...', status: 'pending' },
  { label: 'Generating fixes and advisory report...', status: 'pending' },
])

// API result state
const checkResult = ref<{
  action: string
  compliance_score: number
  risk_category: string
  category: string
  category_state: string
  detected_language: string
  visual_check_confidence: string
  lp_analysis_status: string
  violations: Array<{
    code: string
    title: string
    severity: string
    explanation: string
    suggested_fix: string
    source: string
  }>
  alignment_violations: Array<{
    code: string
    severity: string
    explanation: string
    source_a: string
    source_b: string
    claim_in_source_a: string
    claim_in_source_b: string
  }>
  advisories: Array<{
    code: string
    message: string
  }>
  suggested_copy: string
} | null>(null)
const checkError = ref<string | null>(null)


const adFormat = ref<'single' | 'carousel'>('single')
const selectedPlatform = ref<'Meta' | ''>('Meta')
const selectedRegion = ref('global')
const selectedMinAge = ref('25')
const selectedMaxAge = ref('34')

// UX state
const showLibraryWarning = ref(false)
const isShaking = ref(false)

const isAgeRangeInvalid = computed(() => {
  const min = parseInt(selectedMinAge.value)
  const max = selectedMaxAge.value === '55+' ? 99 : parseInt(selectedMaxAge.value)
  return min > max
})

const parsedDomain = computed(() => {
  if (!landingPageUrl.value) return ''
  try {
    return new URL(landingPageUrl.value).hostname.replace(/^www\./, '').toUpperCase()
  } catch {
    return ''
  }
})

function triggerLibraryAlert() {
  showLibraryWarning.value = true
  setTimeout(() => showLibraryWarning.value = false, 3000)
}

function shakeButton() {
  isShaking.value = true
  setTimeout(() => isShaking.value = false, 500)
}

const adCopy = ref({
  primaryText: '',
  headline: '',
  description: ''
})
const landingPageUrl = ref('')
const urlError = ref('')
const uploadedFile = ref<string | null>(null)
const uploadedFileBlob = ref<File | null>(null)

interface AdCard { visual: string | null; headline: string; description: string }
const carouselCards = ref<AdCard[]>([{ visual: null, headline: '', description: '' }])
const activeCardIndex = ref(0)

function triggerUploadSingle() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*,video/*'
  input.onchange = (e: any) => {
    if (e.target.files && e.target.files[0]) {
      uploadedFileBlob.value = e.target.files[0]
      uploadedFile.value = URL.createObjectURL(e.target.files[0])
    }
  }
  input.click()
}

function triggerUploadCard(index: number) {
  carouselCards.value[index].visual = `https://via.placeholder.com/300x300/1a1a2e/e94560?text=Card+${index + 1}`
}

function addCarouselCard() {
  if (carouselCards.value.length < 10) {
    carouselCards.value.push({ visual: null, headline: '', description: '' })
    activeCardIndex.value = carouselCards.value.length - 1
  }
}

async function nextStep() {
  if (currentStep.value === 1) {
    if (!selectedPlatform.value || isAgeRangeInvalid.value) {
      shakeButton()
      return
    }
    currentStep.value++
    return
  }

  if (currentStep.value === 2) {
    if ((adFormat.value === 'single' && !uploadedFile.value) || (adFormat.value === 'carousel' && !carouselCards.value[0].visual)) {
      shakeButton()
      return
    }
    isChecking.value = true
    currentStep.value = 3
    checkError.value = null
    
    // Reset timer and steps
    timeElapsed.value = 0
    loadingSteps.value.forEach((s, idx) => {
        s.status = idx === 0 ? 'done' : idx === 1 ? 'loading' : 'pending'
    })
    
    timerInterval = setInterval(() => {
        timeElapsed.value += 1
        // Simulate progress based on time
        if (timeElapsed.value === 3) {
            loadingSteps.value[1].status = 'done'
            loadingSteps.value[2].status = 'loading'
        }
        if (timeElapsed.value === 6) {
            loadingSteps.value[2].status = 'done'
            loadingSteps.value[3].status = 'loading'
        }
        if (timeElapsed.value === 11) {
            loadingSteps.value[3].status = 'done'
            loadingSteps.value[4].status = 'loading'
        }
        if (timeElapsed.value === 16) {
            loadingSteps.value[4].status = 'done'
            loadingSteps.value[5].status = 'loading'
        }
    }, 1000)

    try {
      const formData = new FormData()
      formData.append('platform', selectedPlatform.value)
      formData.append('region', selectedRegion.value)
      formData.append('age_min', selectedMinAge.value)
      formData.append('age_max', selectedMaxAge.value)
      formData.append('ad_format', adFormat.value)
      formData.append('primary_text', adCopy.value.primaryText)
      formData.append('headline', adCopy.value.headline)
      formData.append('description', adCopy.value.description)
      formData.append('landing_page_url', landingPageUrl.value)
      
      if (uploadedFileBlob.value) {
        if (uploadedFileBlob.value.type.startsWith('video/')) {
          formData.append('video', uploadedFileBlob.value)
        } else {
          formData.append('image', uploadedFileBlob.value)
        }
      }

      const response = await fetch('http://localhost:8080/checks', {
        method: 'POST',
        // Important: DO NOT set Content-Type header. 
        // The browser automatically sets multipart/form-data and the boundary.
        body: formData
      })

      if (!response.ok) {
        throw new Error('API request failed')
      }

      checkResult.value = await response.json()
      isChecking.value = false
      if (timerInterval) clearInterval(timerInterval)
      currentStep.value = 4
    } catch (err) {
      isChecking.value = false
      if (timerInterval) clearInterval(timerInterval)
      checkError.value = 'Analysis failed. Please check the backend is running and try again.'
      currentStep.value = 2
    }
  } else {
    currentStep.value++
  }
}
</script>

<template>
  <div class="space-y-4">

    <!-- STEP 1: Platform Selection -->
    <div class="border border-border rounded-xl bg-card overflow-hidden">
      <button
        type="button"
        :aria-expanded="currentStep === 1"
        aria-controls="step1-panel"
        @click="currentStep = 1"
        class="w-full flex items-center justify-between p-4 sm:p-6 cursor-pointer hover:bg-muted/30 transition-colors select-none text-left"
      >
        <div class="flex items-center gap-3 sm:gap-4 min-w-0">
          <div class="h-8 w-8 shrink-0 rounded-lg bg-primary text-primary-foreground flex items-center justify-center font-semibold text-sm">
            1
          </div>
          <div class="min-w-0">
            <h2 class="text-base sm:text-lg font-semibold">Select Platform</h2>
            <p class="text-xs sm:text-sm text-muted-foreground hidden sm:block">Choose the platform you want to verify compliance for.</p>
          </div>
        </div>
        <IconChevronDown
          class="text-muted-foreground transition-transform duration-300 shrink-0 ml-2"
          :class="{ 'rotate-180': currentStep === 1 }"
        />
      </button>

      <!-- Accordion Body -->
      <div v-if="currentStep === 1" class="px-6 pb-6 pt-2 border-t border-border animate-in fade-in slide-in-from-top-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Meta (Active/Selectable) -->
          <label 
            class="border-2 rounded-xl p-6 transition-all flex flex-col items-center justify-center text-center gap-3 relative cursor-pointer w-full focus-within:ring-4 focus-within:ring-primary/20 hover:shadow-md active:scale-[0.98]"
            :class="selectedPlatform === 'Meta' ? 'border-primary bg-primary/5 hover:bg-primary/10 shadow-sm' : 'border-border bg-card hover:border-primary/30 hover:bg-muted/30 shadow-sm'"
          >
            <input type="radio" value="Meta" v-model="selectedPlatform" class="absolute opacity-0 w-0 h-0" />
            <div v-if="selectedPlatform === 'Meta'" class="absolute top-3 right-3 text-primary animate-in zoom-in">
              <IconCheck :size="20" stroke-width="3" />
            </div>
            <IconBrandMeta :size="48" :class="selectedPlatform === 'Meta' ? 'text-primary' : 'text-foreground'" />
            <div>
              <h3 class="font-bold text-foreground">Meta</h3>
              <p class="text-xs text-muted-foreground mt-1">Validates compliance within Meta's network (Facebook & Instagram).</p>
            </div>
          </label>

          <!-- Google (Coming Soon) -->
          <div
            class="border border-border rounded-xl p-5 sm:p-6 bg-muted/20 opacity-50 cursor-not-allowed flex flex-col items-center text-center gap-3 relative"
            title="Google Ads support coming in Q3 2026"
          >
            <IconBrandGoogle :size="48" class="text-muted-foreground" />
            <div>
              <h3 class="font-bold text-muted-foreground">Google Ads</h3>
              <p class="text-xs text-muted-foreground mt-1">Coming Soon — Q3 2026</p>
            </div>
          </div>

          <!-- TikTok (Coming Soon) -->
          <div
            class="border border-border rounded-xl p-5 sm:p-6 bg-muted/20 opacity-50 cursor-not-allowed flex flex-col items-center text-center gap-3 relative"
            title="TikTok support coming in Q3 2026"
          >
            <IconBrandTiktok :size="48" class="text-muted-foreground" />
            <div>
              <h3 class="font-bold text-muted-foreground">TikTok</h3>
              <p class="text-xs text-muted-foreground mt-1">Coming Soon — Q3 2026</p>
            </div>
          </div>
        </div>

        <!-- Target Demographics & Region -->
        <div class="mt-8 pt-6 border-t border-border grid grid-cols-1 lg:grid-cols-2 gap-6 pb-2">
           <div class="flex items-center gap-4 bg-muted/30 rounded-xl p-4 border border-border">
             <div class="font-bold text-sm text-foreground whitespace-nowrap min-w-[200px]">Advertising Target Region or Country</div>
             <div class="flex-1 relative">
               <select v-model="selectedRegion" class="w-full appearance-none bg-card border border-border rounded-lg px-4 py-2 pr-10 text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-primary cursor-pointer w-full text-foreground">
                 <option value="global">All Regions & Countries</option>
                 <option value="us">United States</option>
                 <option value="uk">United Kingdom</option>
                 <option value="eu">European Union</option>
               </select>
               <IconChevronDown :size="16" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none" />
             </div>
           </div>
           
           <div class="flex items-center gap-4 bg-muted/30 rounded-xl p-4 border border-border">
             <div class="flex items-center gap-2 font-bold text-sm text-foreground whitespace-nowrap">
               <IconUsers :size="20" class="text-muted-foreground" />
               Age Range
             </div>
             <div class="flex-1 flex items-center gap-2 ml-4">
               <div class="relative flex-1">
                 <select v-model="selectedMinAge" class="w-full appearance-none bg-card border border-border rounded-lg px-4 py-2 pr-10 text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-primary cursor-pointer text-foreground">
                   <option value="18">18</option>
                   <option value="25">25</option>
                   <option value="35">35</option>
                 </select>
                 <IconChevronDown :size="16" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none" />
               </div>
               <span class="text-muted-foreground font-medium">-</span>
               <div class="relative flex-1">
                 <select v-model="selectedMaxAge" class="w-full appearance-none bg-card border border-border rounded-lg px-4 py-2 pr-10 text-sm font-semibold focus:outline-none focus:ring-2 focus:ring-primary cursor-pointer text-foreground">
                   <option value="24">24</option>
                   <option value="34">34</option>
                   <option value="44">44</option>
                   <option value="55+">55+</option>
                 </select>
                 <IconChevronDown :size="16" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground pointer-events-none" />
               </div>
             </div>
              </div>
          </div>

        <div v-if="isAgeRangeInvalid" class="mt-4 p-3 bg-destructive/10 border border-destructive/20 rounded-lg flex items-center gap-2 animate-in fade-in slide-in-from-top-2">
          <IconAlertCircle class="text-destructive" :size="16" />
          <span class="text-xs font-bold text-destructive">Invalid age range: Minimum age cannot be greater than maximum age.</span>
        </div>

        <div class="mt-6 flex justify-end">
          <button 
            @click="nextStep" 
            :disabled="!selectedPlatform" 
            class="bg-primary text-primary-foreground px-6 py-2.5 rounded-lg font-medium hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            :class="{ 'animate-shake bg-destructive': isShaking }"
          >
            Next Step
          </button>
        </div>
      </div>
    </div>

    <!-- STEP 2: Creative & Ad Copy -->
    <div class="border border-border rounded-xl bg-card overflow-hidden">
      <button
        type="button"
        :aria-expanded="currentStep === 2"
        aria-controls="step2-panel"
        @click="currentStep = 2"
        class="w-full flex items-center justify-between p-4 sm:p-6 cursor-pointer hover:bg-muted/30 transition-colors text-left"
      >
        <div class="flex items-center gap-3 sm:gap-4 min-w-0">
          <div
            class="h-8 w-8 shrink-0 rounded-lg flex items-center justify-center font-semibold text-sm transition-colors"
            :class="currentStep >= 2 ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'"
          >
            2
          </div>
          <div class="min-w-0">
            <h2 class="text-base sm:text-lg font-semibold" :class="currentStep < 2 && 'text-muted-foreground'">Ad Creative & Copy</h2>
            <p class="text-xs sm:text-sm text-muted-foreground hidden sm:block">Upload visual and provide the exact ad text for full compliance check.</p>
          </div>
        </div>
        <IconChevronDown
          class="text-muted-foreground transition-transform duration-300 shrink-0 ml-2"
          :class="{ 'rotate-180': currentStep === 2 }"
        />
      </button>

      <div
        v-if="currentStep === 2"
        id="step2-panel"
        role="region"
        class="px-4 sm:px-6 pb-6 pt-2 border-t border-border animate-in fade-in slide-in-from-top-4"
      >
        <!-- Ad Format Selector -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8 pb-4 border-b border-border">
          <div>
            <h3 class="font-bold text-lg text-foreground">Ad Format & Creative</h3>
            <p class="text-sm text-muted-foreground mt-1">Select your ad format and configure the visuals and copy.</p>
          </div>
          <div class="flex bg-muted/30 p-1 rounded-xl border border-border shrink-0">
            <button
              type="button"
              @click="adFormat = 'single'"
              :class="{'bg-card text-foreground shadow-sm ring-1 ring-border': adFormat === 'single', 'text-muted-foreground hover:text-foreground': adFormat !== 'single'}"
              class="px-4 py-2 rounded-lg text-sm font-bold transition-all"
            >
              Single Image or Video
            </button>
            <button
              type="button"
              @click="adFormat = 'carousel'"
              :class="{'bg-card text-foreground shadow-sm ring-1 ring-border': adFormat === 'carousel', 'text-muted-foreground hover:text-foreground': adFormat !== 'carousel'}"
              class="px-4 py-2 rounded-lg text-sm font-bold transition-all"
            >
              Carousel
            </button>
          </div>
        </div>

        <!-- SINGLE IMAGE FORMAT UI -->
        <div v-if="adFormat === 'single'" class="grid grid-cols-1 lg:grid-cols-2 gap-6 sm:gap-8">
          <div class="space-y-3">
            <h3 class="font-bold text-sm text-foreground">Visual Creative</h3>

            <div
              @click="triggerUploadSingle"
              class="border-2 border-dashed border-border rounded-xl p-8 hover:bg-muted/50 transition-colors cursor-pointer flex flex-col items-center justify-center text-center group"
              :class="uploadedFile ? 'bg-muted/30 border-primary/50 min-h-[200px]' : 'min-h-[280px]'"
              role="button"
              tabindex="0"
              aria-label="Click to upload image or video"
              @keydown.enter="triggerUploadSingle"
              @keydown.space.prevent="triggerUploadSingle"
            >
              <template v-if="!uploadedFile">
                <div class="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center text-primary mb-4 group-hover:scale-110 transition-transform">
                  <IconUpload />
                </div>
                <p class="font-bold text-foreground">Upload your image or video</p>
                <p class="text-sm text-muted-foreground mt-1">Click here or drag & drop</p>
                
                <div class="relative mt-6">
                  <button 
                    @click.stop="triggerLibraryAlert" 
                    type="button" 
                    class="bg-secondary text-secondary-foreground px-5 py-2.5 rounded-lg text-sm font-bold hover:bg-secondary/80 transition-colors"
                  >
                    Image Library
                  </button>
                  <div v-if="showLibraryWarning" class="absolute top-full mt-2 left-1/2 -translate-x-1/2 whitespace-nowrap bg-foreground text-background text-[10px] py-1 px-2 rounded shadow-lg animate-in fade-in zoom-in">
                    Gallery Coming Soon
                  </div>
                </div>
              </template>
              <template v-else>
                <div class="w-full h-full flex flex-col items-center justify-center space-y-4">
                  <div class="w-56 h-56 bg-primary/20 rounded-xl flex items-center justify-center border-2 border-primary/30 text-primary font-bold shadow-inner relative overflow-hidden">
                    <img :src="uploadedFile" class="object-cover w-full h-full opacity-80" />
                    <div class="absolute inset-0 flex items-center justify-center bg-black/40"><span class="bg-background px-3 py-1 rounded text-xs truncate max-w-[80%]">{{ uploadedFileBlob?.name }}</span></div>
                  </div>
                  <span class="text-sm font-bold text-primary cursor-pointer hover:underline">Change File</span>
                </div>
              </template>
            </div>
          </div>

          <div class="space-y-4 bg-muted/20 p-5 sm:p-6 border border-border rounded-xl shadow-sm">
            <div class="border-b border-border pb-3 mb-2 flex items-center gap-2">
              <IconBrandMeta class="text-primary" />
              <h3 class="font-bold text-foreground">Meta Ad Copy</h3>
            </div>

            <div class="space-y-1.5">
              <label for="primary-text" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Primary Text</label>
              <textarea
                id="primary-text"
                v-model="adCopy.primaryText"
                rows="4"
                class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                placeholder="Appears above your ad..."
              ></textarea>
            </div>
            <div class="space-y-1.5">
              <label for="headline" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Headline</label>
              <input
                id="headline"
                v-model="adCopy.headline"
                type="text"
                class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                placeholder="Usually appears below the image..."
              />
            </div>
            <div class="space-y-1.5">
              <label for="description" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Description <span class="normal-case font-normal text-muted-foreground/70">(Optional)</span></label>
              <input
                id="description"
                v-model="adCopy.description"
                type="text"
                class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                placeholder="Additional context..."
              />
            </div>
            <div class="space-y-1.5 pt-2 border-t border-border mt-2">
              <label for="landing-url" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">
                Landing Page URL <span class="normal-case font-normal text-muted-foreground/70">(Recommended)</span>
              </label>
              <input
                id="landing-url"
                v-model="landingPageUrl"
                type="url"
                class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                :class="urlError ? 'border-destructive focus:ring-destructive' : ''"
                placeholder="https://yoursite.com/product-page"
                @blur="validateUrl"
              />
              <p v-if="urlError" class="text-xs text-destructive font-medium">{{ urlError }}</p>
              <p v-else class="text-[11px] text-muted-foreground">Meta checks your landing page too. Skipping this may miss destination-level violations.</p>
            </div>
          </div>
        </div>

        <!-- CAROUSEL FORMAT UI -->
        <div v-else class="space-y-8">
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label for="carousel-primary-text" class="text-xs font-bold text-muted-foreground uppercase tracking-wider flex items-center gap-2">
                <IconBrandMeta class="text-primary w-4 h-4" /> Global Primary Text
              </label>
              <textarea
                id="carousel-primary-text"
                v-model="adCopy.primaryText"
                rows="2"
                class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                placeholder="This text appears above the entire carousel..."
              ></textarea>
            </div>
            <div class="space-y-1.5">
              <label for="carousel-landing-url" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">
                Landing Page URL <span class="normal-case font-normal text-muted-foreground/70">(Recommended)</span>
              </label>
              <input
                id="carousel-landing-url"
                v-model="landingPageUrl"
                type="url"
                class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                :class="urlError ? 'border-destructive focus:ring-destructive' : ''"
                placeholder="https://yoursite.com/product-page"
                @blur="validateUrl"
              />
              <p v-if="urlError" class="text-xs text-destructive font-medium">{{ urlError }}</p>
              <p v-else class="text-[11px] text-muted-foreground">Meta checks your landing page too. Skipping this may miss destination-level violations.</p>
            </div>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 sm:gap-8">
            <div class="lg:col-span-2 border border-border rounded-xl bg-muted/10 p-5 overflow-hidden flex flex-col">
              <h3 class="font-bold text-sm text-foreground mb-4">Carousel Cards ({{ carouselCards.length }}/10)</h3>
              <div class="flex gap-4 overflow-x-auto pb-4 flex-1 items-center">
                <div
                  v-for="(card, idx) in carouselCards"
                  :key="idx"
                  @click="activeCardIndex = idx"
                  class="shrink-0 w-44 h-60 rounded-xl border-2 transition-all cursor-pointer flex flex-col items-center justify-center relative bg-card overflow-hidden group"
                  :class="activeCardIndex === idx ? 'border-primary shadow-md scale-105' : 'border-dashed border-border hover:border-solid hover:border-primary/50'"
                >
                  <div v-if="activeCardIndex === idx" class="absolute top-2 right-2 bg-primary text-primary-foreground text-[10px] font-bold px-2 py-0.5 rounded-full z-10">Editing</div>
                  <template v-if="!card.visual">
                    <IconUpload class="text-muted-foreground group-hover:text-primary transition-colors mb-2" />
                    <span class="text-xs font-bold text-muted-foreground">Click to upload</span>
                    <span class="text-[10px] text-muted-foreground/70 mt-1">Card {{ idx + 1 }}</span>
                  </template>
                  <template v-else>
                    <img :src="card.visual" class="w-full h-full object-cover opacity-90" :alt="`Card ${idx + 1}`" />
                    <div class="absolute bottom-0 w-full bg-background/90 backdrop-blur border-t border-border p-2">
                      <div class="text-[10px] font-bold truncate">{{ card.headline || 'Headline' }}</div>
                    </div>
                  </template>
                </div>

                <div
                  v-if="carouselCards.length < 10"
                  @click="addCarouselCard"
                  class="shrink-0 w-16 h-60 border-2 border-dashed border-border rounded-xl hover:bg-muted/50 transition-colors cursor-pointer flex items-center justify-center text-muted-foreground hover:text-foreground"
                  role="button"
                  aria-label="Add carousel card"
                >
                  <span class="text-2xl font-light">+</span>
                </div>
              </div>
            </div>

            <div class="space-y-4 bg-muted/20 p-5 border border-border rounded-xl shadow-sm h-fit">
              <div class="border-b border-border pb-3 mb-2 flex items-center justify-between">
                <h3 class="font-bold text-foreground text-sm">Editing Card {{ activeCardIndex + 1 }}</h3>
              </div>
              <div class="space-y-1.5">
                <label :for="`card-media-${activeCardIndex}`" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Card Media</label>
                <div
                  :id="`card-media-${activeCardIndex}`"
                  @click="triggerUploadCard(activeCardIndex)"
                  class="w-full h-24 border border-dashed border-border rounded-lg bg-card hover:bg-muted/50 cursor-pointer flex items-center justify-center transition-colors"
                  role="button"
                  :aria-label="`Upload media for card ${activeCardIndex + 1}`"
                >
                  <span v-if="!carouselCards[activeCardIndex].visual" class="text-xs font-bold text-muted-foreground flex items-center gap-2"><IconUpload :size="14" /> Upload Media</span>
                  <span v-else class="text-xs font-bold text-primary flex items-center gap-2">Replace Image</span>
                </div>
              </div>
              <div class="space-y-1.5">
                <label :for="`card-headline-${activeCardIndex}`" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Card Headline</label>
                <input
                  :id="`card-headline-${activeCardIndex}`"
                  v-model="carouselCards[activeCardIndex].headline"
                  type="text"
                  class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                  placeholder="Short headline..."
                />
              </div>
              <div class="space-y-1.5">
                <label :for="`card-description-${activeCardIndex}`" class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Description</label>
                <input
                  :id="`card-description-${activeCardIndex}`"
                  v-model="carouselCards[activeCardIndex].description"
                  type="text"
                  class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow"
                  placeholder="Sub-text..."
                />
              </div>
            </div>
          </div>
        </div>

        <div class="mt-8 flex justify-end pt-6 border-t border-border">
          <button 
            @click="nextStep" 
            :disabled="(adFormat === 'single' ? !uploadedFile : !carouselCards[0].visual)" 
            class="bg-[#ed3875] text-white px-8 py-3 rounded-lg font-bold hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-md"
            :class="{ 'animate-shake bg-destructive': isShaking }"
          >
            Check Compliance >
          </button>
        </div>
      </div>
    </div>

    <!-- STEP 3: Processing State -->
    <div v-if="currentStep === 3" class="border border-border rounded-xl bg-card p-12 flex flex-col items-center justify-center text-center animate-in fade-in zoom-in-95">
       <IconLoader2 class="animate-spin text-primary mb-6" :size="48" />
       <h2 class="text-xl font-bold text-primary mb-2">AI is analyzing your ad...</h2>
       <p class="text-muted-foreground tabular-nums">Elapsed Time: <span class="font-bold text-foreground">{{ timeElapsed }}s</span></p>
       
       <div class="mt-8 space-y-3 w-full max-w-md">
         <div v-for="(step, index) in loadingSteps" :key="index"
              class="bg-card border border-border p-3 rounded-lg flex items-center justify-between shadow-sm transition-all duration-300"
              :class="{ 'opacity-50': step.status === 'pending' }">
           <span class="text-sm text-foreground" :class="{'font-bold text-primary': step.status === 'loading'}">{{ step.label }}</span>
           <IconCheck v-if="step.status === 'done'" class="text-green-600" :size="16" />
           <IconLoader2 v-else-if="step.status === 'loading'" class="text-primary animate-spin" :size="16" />
           <div v-else class="w-4 h-4 rounded-full border-2 border-muted" />
         </div>
       </div>
    </div>

    <!-- STEP 4: Results Dashboard -->
    <div v-if="currentStep === 4" class="animate-in fade-in slide-in-from-bottom-4">
      <div class="flex items-center justify-between mb-8 pb-4 border-b border-border">
        <div>
          <h2 class="text-2xl font-bold text-foreground">Compliance Report</h2>
          <p class="text-sm text-muted-foreground mt-1 tracking-tight">Review the issues found within your ad creative and copy.</p>
        </div>
        <button 
          @click="currentStep = 1" 
          class="bg-muted hover:bg-muted/80 text-foreground px-4 py-2 rounded-lg text-sm font-bold border border-border transition-colors flex items-center gap-2"
        >
          Start New Check
        </button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
        
        <!-- Left: Context/Preview (2 cols) -->
        <div class="col-span-2">
          <div class="sticky top-8 space-y-4">
            <div class="border border-border rounded-xl bg-card p-4 shadow-sm relative overflow-hidden">
              <div class="absolute top-0 left-0 w-1 h-full bg-primary/20"></div>
            <!-- Simulated Facebook Ad Card -->
            <div class="flex items-center gap-2 mb-3">
              <div class="w-8 h-8 rounded-full bg-muted shrink-0"></div>
              <div>
                <div class="text-[10px] text-muted-foreground">Sponsored</div>
              </div>
            </div>
            <p v-if="adCopy.primaryText" class="text-sm text-foreground mb-3">{{ adCopy.primaryText }}</p>
            <div class="w-full max-h-64 aspect-video bg-muted/30 rounded overflow-hidden flex items-center justify-center shadow-inner border border-border">
              <img v-if="uploadedFile" :src="uploadedFile" class="w-full h-full object-cover" />
              <span v-else class="text-xs text-muted-foreground font-medium">No image uploaded</span>
            </div>
            <div class="bg-muted/30 p-3 mt-2 border border-border flex justify-between items-center gap-2">
              <div class="min-w-0">
                <p v-if="parsedDomain" class="text-xs text-muted-foreground uppercase truncate">{{ parsedDomain }}</p>
                <p v-if="adCopy.headline" class="text-sm font-bold truncate">{{ adCopy.headline }}</p>
                <p v-if="adCopy.description" class="text-xs text-muted-foreground truncate">{{ adCopy.description }}</p>
              </div>
               <button class="bg-secondary text-secondary-foreground text-xs font-bold px-3 py-1 rounded shrink-0">Learn More</button>
            </div>
          </div>
        </div>
        </div>

        <!-- Right: Results Checklist (3 cols) -->
        <div class="col-span-3 space-y-4">

          <!-- Score Summary Card -->
          <div class="border border-border rounded-xl bg-card p-5">
            <div class="flex items-center justify-between mb-3">
              <div>
                <div class="text-3xl font-bold" :class="checkResult?.compliance_score === 0 ? 'text-destructive' : checkResult?.compliance_score >= 70 ? 'text-green-600' : 'text-amber-500'">
                  {{ checkResult?.compliance_score ?? 0 }}<span class="text-lg font-normal text-muted-foreground">/100</span>
                </div>
                <div class="text-sm font-semibold mt-1" :class="checkResult?.risk_category === 'High Risk' ? 'text-destructive' : checkResult?.risk_category === 'Medium Risk' ? 'text-amber-500' : 'text-green-600'">
                  {{ checkResult?.risk_category }}
                </div>
              </div>
              <div class="text-right">
                <div class="text-xs text-muted-foreground">Category</div>
                <div class="text-sm font-semibold">{{ checkResult?.category }}</div>
                <div class="text-xs mt-1 px-2 py-0.5 rounded-full inline-block"
                  :class="checkResult?.action === 'PROHIBITED' ? 'bg-destructive/10 text-destructive' : checkResult?.action === 'AUTHORIZATION_REQUIRED' ? 'bg-orange-100 text-orange-700' : checkResult?.action === 'RESTRICTED' ? 'bg-amber-100 text-amber-700' : checkResult?.action === 'ALLOWED_WITH_RESTRICTIONS' ? 'bg-blue-100 text-blue-700' : 'bg-green-100 text-green-700'">
                  {{ checkResult?.action }}
                </div>
              </div>
            </div>
            <!-- Score bar -->
            <div class="w-full bg-muted rounded-full h-2">
              <div class="h-2 rounded-full transition-all duration-500"
                :class="checkResult?.compliance_score === 0 ? 'bg-destructive' : checkResult?.compliance_score >= 70 ? 'bg-green-500' : 'bg-amber-400'"
                :style="`width: ${checkResult?.compliance_score ?? 0}%`">
              </div>
            </div>
          </div>

          <AdvisoryPanel v-if="checkResult?.advisories" :advisories="checkResult.advisories" />

          <div v-if="checkResult?.suggested_copy" class="border border-green-200 rounded-xl bg-green-50 p-5">
            <h3 class="font-bold text-green-800 text-sm mb-2">Suggested Compliant Copy</h3>
            <p class="text-sm text-green-900 bg-white/50 p-3 rounded-lg border border-green-100 whitespace-pre-line">{{ checkResult.suggested_copy }}</p>
          </div>
          
          <!-- Alignment Violations List -->
          <div v-if="checkResult?.alignment_violations && checkResult.alignment_violations.length > 0"
               class="border border-amber-200 rounded-xl bg-card overflow-hidden">
            <div class="bg-amber-50 border-b border-amber-200 p-4 flex items-center gap-2">
              <IconAlertCircle :size="18" class="text-amber-600 shrink-0" />
              <span class="font-bold text-amber-800 text-sm">Cross-Source Contradictions (Alignment)</span>
            </div>
            <div class="p-6 space-y-6">
              <div v-for="violation in checkResult.alignment_violations" :key="'align'+violation.code">
                <div class="flex items-start gap-2">
                  <div class="flex-1">
                    <h4 class="font-bold text-sm text-foreground mb-1">{{ violation.explanation }}</h4>
                    <div class="grid grid-cols-2 gap-4 mt-3">
                      <div class="bg-muted p-3 rounded border border-border">
                        <span class="text-[10px] uppercase font-bold text-muted-foreground block mb-1">In {{ violation.source_a }}</span>
                        <p class="text-xs text-foreground">{{ violation.claim_in_source_a }}</p>
                      </div>
                      <div class="bg-muted p-3 rounded border border-border">
                        <span class="text-[10px] uppercase font-bold text-muted-foreground block mb-1">In {{ violation.source_b }}</span>
                        <p class="text-xs text-foreground">{{ violation.claim_in_source_b }}</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Violations List -->
          <div v-if="checkResult?.violations && checkResult.violations.length > 0"
               class="border border-border rounded-xl bg-card overflow-hidden">
            <div class="bg-destructive/10 border-b border-destructive/20 p-4 flex items-center justify-between">
              <div class="flex items-center gap-2 text-destructive font-bold">
                <IconAlertCircle :size="20" />
                {{ checkResult.violations.length }} Issue{{ checkResult.violations.length > 1 ? 's' : '' }} Found
              </div>
              <div class="bg-destructive text-destructive-foreground text-xs font-bold px-2 py-0.5 rounded-full">
                {{ checkResult.risk_category }}
              </div>
            </div>
            <div class="p-6 space-y-6">
              <div v-for="violation in checkResult.violations" :key="violation.code">
                <div class="flex items-start gap-2">
                  <IconAlertCircle class="text-destructive mt-0.5 shrink-0" :size="16" />
                  <div class="flex-1">
                    <div class="flex items-center gap-2 mb-1">
                      <h4 class="font-bold text-sm text-foreground">{{ violation.title }}</h4>
                      <span class="text-[10px] font-bold px-1.5 py-0.5 rounded"
                        :class="violation.severity === 'PROHIBITED' ? 'bg-destructive/10 text-destructive' : violation.severity === 'AUTHORIZATION_REQUIRED' ? 'bg-orange-100 text-orange-700' : violation.severity === 'RESTRICTED' ? 'bg-amber-100 text-amber-700' : 'bg-blue-100 text-blue-700'">
                        {{ violation.severity }}
                      </span>
                    </div>
                    <p class="text-sm text-muted-foreground">{{ violation.explanation }}</p>
                    
                    <!-- Source Badge -->
                    <div class="mt-2 text-[10px] font-bold text-muted-foreground uppercase flex items-center gap-1">
                      <IconAlertCircle :size="12" /> Found in: 
                      <span class="bg-muted px-1.5 py-0.5 rounded">{{ violation.source }}</span>
                    </div>

                    <div v-if="violation.suggested_fix" class="mt-2 bg-green-50 border border-green-200 rounded-lg p-3">
                      <p class="text-xs font-bold text-green-700 mb-1">Suggested Fix</p>
                      <p class="text-sm text-green-800">{{ violation.suggested_fix }}</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- No violations -->
          <div v-else class="bg-green-50 border border-green-200 p-4 rounded-xl flex items-start gap-3">
            <IconCheck class="text-green-600 shrink-0 mt-0.5" />
            <p class="text-sm text-green-800 font-medium">No policy violations detected. This ad appears compliant with Meta's advertising standards.</p>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-4px); }
  75% { transform: translateX(4px); }
}

.animate-shake {
  animation: shake 0.2s ease-in-out 0s 2;
}

/* Ensure sticky preview doesn't overflow container on very short viewports */
.sticky {
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  scrollbar-width: none;
}
.sticky::-webkit-scrollbar {
  display: none;
}
</style>

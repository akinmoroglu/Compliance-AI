<script setup lang="ts">
import { ref } from 'vue'
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

// Critically Restored State Variables
const currentStep = ref(1)
const isChecking = ref(false)

const adFormat = ref<'single' | 'carousel'>('single')
const selectedPlatform = ref<'Meta' | ''>('')
const selectedRegion = ref('global')
const selectedMinAge = ref('25')
const selectedMaxAge = ref('34')

const adCopy = ref({
  primaryText: '',
  headline: '',
  description: ''
})
const uploadedFile = ref<string | null>(null)

interface AdCard {
  visual: string | null;
  headline: string;
  description: string;
}
const carouselCards = ref<AdCard[]>([
  { visual: null, headline: '', description: '' }
])
const activeCardIndex = ref(0)

function triggerUploadSingle() {
  uploadedFile.value = 'mock-creative.jpg'
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

function nextStep() {
  if (currentStep.value === 2) {
    // Simulate checking
    isChecking.value = true
    currentStep.value = 3
    setTimeout(() => {
      isChecking.value = false
      currentStep.value = 4
    }, 2500)
  } else {
    currentStep.value++
  }
}
</script>

<template>
  <div class="space-y-4">
    <!-- STEP 1: Platform Selection -->
    <div class="border border-border rounded-xl bg-card overflow-hidden">
      <!-- Accordion Header -->
      <button 
        type="button"
        @click="currentStep = 1"
        class="w-full flex items-center justify-between p-6 cursor-pointer hover:bg-muted/30 transition-colors select-none text-left"
      >
        <div class="flex items-center gap-4">
          <div class="h-8 w-8 rounded-lg bg-primary text-primary-foreground flex items-center justify-center font-semibold text-sm">
            1
          </div>
          <div>
            <h2 class="text-lg font-semibold cursor-pointer">Select Platform</h2>
            <p class="text-sm text-muted-foreground">Choose the platform you want to verify compliance for.</p>
          </div>
        </div>
        <IconChevronDown 
          class="text-muted-foreground transition-transform duration-300"
          :class="{ 'rotate-180': currentStep === 1 }"
        />
      </button>

      <!-- Accordion Body -->
      <div v-if="currentStep === 1" class="px-6 pb-6 pt-2 border-t border-border animate-in fade-in slide-in-from-top-4">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Meta (Active/Selectable) -->
          <label 
            class="border-2 rounded-xl p-6 transition-colors flex flex-col items-center justify-center text-center gap-3 relative cursor-pointer w-full focus-within:ring-4 focus-within:ring-primary/20"
            :class="selectedPlatform === 'Meta' ? 'border-primary bg-primary/5 hover:bg-primary/10 shadow-sm' : 'border-border bg-card hover:bg-muted/30 shadow-sm'"
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

          <!-- Google (Disabled) -->
          <div class="border border-border rounded-xl p-6 bg-muted/20 opacity-50 cursor-not-allowed flex flex-col items-center text-center gap-3 relative">
            <IconBrandGoogle :size="48" class="text-muted-foreground" />
            <div>
              <h3 class="font-bold text-muted-foreground">Google Ads</h3>
              <p class="text-xs text-muted-foreground mt-1">Coming Soon.</p>
            </div>
          </div>

          <!-- TikTok (Disabled) -->
          <div class="border border-border rounded-xl p-6 bg-muted/20 opacity-50 cursor-not-allowed flex flex-col items-center text-center gap-3 relative">
            <IconBrandTiktok :size="48" class="text-muted-foreground" />
            <div>
              <h3 class="font-bold text-muted-foreground">TikTok</h3>
              <p class="text-xs text-muted-foreground mt-1">Coming Soon.</p>
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

        <div class="mt-6 flex justify-end">
          <button @click="nextStep" :disabled="!selectedPlatform" class="bg-primary text-primary-foreground px-6 py-2.5 rounded-lg font-medium hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed">
            Next Step
          </button>
        </div>
      </div>
    </div>

    <!-- STEP 2: Creative & Ad Copy -->
    <div class="border border-border rounded-xl bg-card overflow-hidden">
      <!-- Accordion Header -->
      <div 
        @click="currentStep = 2"
        class="flex items-center justify-between p-6 cursor-pointer hover:bg-muted/30 transition-colors"
      >
        <div class="flex items-center gap-4">
          <div 
            class="h-8 w-8 rounded-lg flex items-center justify-center font-semibold text-sm transition-colors"
            :class="currentStep >= 2 ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'"
          >
            2
          </div>
          <div>
            <h2 class="text-lg font-semibold" :class="currentStep < 2 && 'text-muted-foreground'">Ad Creative & Copy</h2>
            <p class="text-sm text-muted-foreground">Upload visual and provide the exact ad text for full compliance check.</p>
          </div>
        </div>
        <IconChevronDown 
          class="text-muted-foreground transition-transform duration-300"
          :class="{ 'rotate-180': currentStep === 2 }"
        />
      </div>

      <!-- Accordion Body -->
      <div v-if="currentStep === 2" class="px-6 pb-6 pt-2 border-t border-border animate-in fade-in slide-in-from-top-4">
        
        <!-- Ad Format Selector -->
        <div class="flex items-center justify-between mb-8 pb-4 border-b border-border">
          <div>
            <h3 class="font-bold text-lg text-foreground">Ad Format & Creative</h3>
            <p class="text-sm text-muted-foreground mt-1">Select your ad format and configure the visuals and copy.</p>
          </div>
          <div class="flex bg-muted/30 p-1 rounded-xl border border-border">
            <button 
              @click="adFormat = 'single'" 
              :class="{'bg-card text-foreground shadow-sm ring-1 ring-border': adFormat === 'single', 'text-muted-foreground hover:text-foreground': adFormat !== 'single'}" 
              class="px-5 py-2 rounded-lg text-sm font-bold transition-all"
            >
              Single Image or Video
            </button>
            <button 
              @click="adFormat = 'carousel'" 
              :class="{'bg-card text-foreground shadow-sm ring-1 ring-border': adFormat === 'carousel', 'text-muted-foreground hover:text-foreground': adFormat !== 'carousel'}" 
              class="px-5 py-2 rounded-lg text-sm font-bold transition-all"
            >
              Carousel Format
            </button>
          </div>
        </div>

        <!-- SINGLE IMAGE FORMAT UI -->
        <div v-if="adFormat === 'single'" class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div class="space-y-3">
            <h3 class="font-bold text-sm text-foreground">Visual Creative</h3>
            <div 
              @click="triggerUploadSingle"
              class="border-2 border-dashed border-border rounded-xl p-8 hover:bg-muted/50 transition-colors cursor-pointer flex flex-col items-center justify-center text-center group min-h-[350px]"
              :class="{ 'bg-muted/30 border-primary/50': uploadedFile }"
            >
              <template v-if="!uploadedFile">
                <div class="h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center text-primary mb-4 group-hover:scale-110 transition-transform">
                  <IconUpload />
                </div>
                <p class="font-bold text-foreground">Upload your image or video</p>
                <p class="text-sm text-muted-foreground mt-1">Click here or drag & drop</p>
                <button class="mt-6 bg-secondary text-secondary-foreground px-5 py-2.5 rounded-lg text-sm font-bold">Image Library</button>
              </template>
              <template v-else>
                <div class="w-full h-full flex flex-col items-center justify-center space-y-4">
                  <div class="w-56 h-56 bg-primary/20 rounded-xl flex items-center justify-center border-2 border-primary/30 text-primary font-bold shadow-inner relative overflow-hidden">
                    <img src="https://via.placeholder.com/300x500/1a1a2e/e94560?text=Mock+Ad" class="object-cover w-full h-full opacity-80" />
                    <div class="absolute inset-0 flex items-center justify-center bg-black/40"><span class="bg-background px-3 py-1 rounded text-xs">mock-creative.jpg</span></div>
                  </div>
                  <span class="text-sm font-bold text-primary cursor-pointer hover:underline">Change File</span>
                </div>
              </template>
            </div>
          </div>

          <div class="space-y-4 bg-muted/20 p-6 border border-border rounded-xl shadow-sm">
            <div class="border-b border-border pb-3 mb-2 flex items-center gap-2">
              <IconBrandMeta class="text-primary" />
              <h3 class="font-bold text-foreground">Meta Ad Copy</h3>
            </div>
            
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Primary Text</label>
              <textarea v-model="adCopy.primaryText" rows="4" class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow" placeholder="Appears above your ad..."></textarea>
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Headline</label>
              <input v-model="adCopy.headline" type="text" class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow" placeholder="Usually appears below the image..." />
            </div>
            <div class="space-y-1.5">
              <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Description (Optional)</label>
              <input v-model="adCopy.description" type="text" class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow" placeholder="Additional context..." />
            </div>
          </div>
        </div>

        <!-- CAROUSEL FORMAT UI -->
        <div v-else class="space-y-8">
          <!-- Global Carousel Copy -->
          <div class="space-y-1.5">
            <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider flex items-center gap-2"><IconBrandMeta class="text-primary w-4 h-4" /> Global Primary Text</label>
            <textarea v-model="adCopy.primaryText" rows="2" class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow" placeholder="This text appears above the entire carousel..."></textarea>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
            <!-- Carousel Cards Scroller (Left 2 cols) -->
            <div class="lg:col-span-2 border border-border rounded-xl bg-muted/10 p-6 overflow-hidden flex flex-col">
              <h3 class="font-bold text-sm text-foreground mb-4">Carousel Cards ({{ carouselCards.length }}/10)</h3>
              <div class="flex gap-4 overflow-x-auto pb-4 flex-1 items-center">
                 
                 <!-- Render actual cards -->
                 <div 
                   v-for="(card, idx) in carouselCards" 
                   :key="idx"
                   @click="activeCardIndex = idx"
                   class="shrink-0 w-44 h-60 rounded-xl border-2 transition-all cursor-pointer flex flex-col items-center justify-center relative bg-card overflow-hidden group"
                   :class="activeCardIndex === idx ? 'border-primary shadow-md scale-105' : 'border-dashed border-border hover:border-solid hover:border-primary/50'"
                 >
                    <!-- Active Indicator -->
                    <div v-if="activeCardIndex === idx" class="absolute top-2 right-2 bg-primary text-primary-foreground text-[10px] font-bold px-2 py-0.5 rounded-full z-10">Editing</div>
                    
                    <template v-if="!card.visual">
                       <IconUpload class="text-muted-foreground group-hover:text-primary transition-colors mb-2" />
                       <span class="text-xs font-bold text-muted-foreground">Click to upload</span>
                       <span class="text-[10px] text-muted-foreground/70 mt-1">Card {{ idx + 1 }}</span>
                    </template>
                    <template v-else>
                       <img :src="card.visual" class="w-full h-full object-cover opacity-90" />
                       <div class="absolute bottom-0 w-full bg-background/90 backdrop-blur border-t border-border p-2">
                         <div class="text-[10px] font-bold truncate">{{ card.headline || 'Headline' }}</div>
                       </div>
                    </template>
                 </div>

                 <!-- Add Card Button -->
                 <div 
                   v-if="carouselCards.length < 10"
                   @click="addCarouselCard"
                   class="shrink-0 w-16 h-60 border-2 border-dashed border-border rounded-xl hover:bg-muted/50 transition-colors cursor-pointer flex items-center justify-center text-muted-foreground hover:text-foreground"
                 >
                   <span class="text-2xl font-light">+</span>
                 </div>
              </div>
            </div>

            <!-- Active Card Editor (Right 1 col) -->
            <div class="space-y-4 bg-muted/20 p-6 border border-border rounded-xl shadow-sm h-fit relative">
               <!-- Editor Header -->
               <div class="border-b border-border pb-3 mb-2 flex items-center justify-between">
                 <h3 class="font-bold text-foreground text-sm">Editing Card {{ activeCardIndex + 1 }}</h3>
               </div>
               
               <!-- Visual Upload action for this specific card -->
               <div class="space-y-1.5">
                  <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Card Media</label>
                  <div 
                    @click="triggerUploadCard(activeCardIndex)"
                    class="w-full h-24 border border-dashed border-border rounded-lg bg-card hover:bg-muted/50 cursor-pointer flex items-center justify-center transition-colors"
                  >
                     <span v-if="!carouselCards[activeCardIndex].visual" class="text-xs font-bold text-muted-foreground flex items-center gap-2"><IconUpload :size="14" /> Upload Media</span>
                     <span v-else class="text-xs font-bold text-primary flex items-center gap-2">Replace Image</span>
                  </div>
               </div>

               <!-- Card specific copy -->
               <div class="space-y-1.5">
                 <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Card Headline</label>
                 <input v-model="carouselCards[activeCardIndex].headline" type="text" class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow" placeholder="Short headline..." />
               </div>
               <div class="space-y-1.5">
                 <label class="text-xs font-bold text-muted-foreground uppercase tracking-wider">Description</label>
                 <input v-model="carouselCards[activeCardIndex].description" type="text" class="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:border-transparent transition-shadow" placeholder="Sub-text..." />
               </div>
            </div>
          </div>
        </div>

        <div class="mt-8 flex justify-end pt-6 border-t border-border">
          <button @click="nextStep" :disabled="(adFormat === 'single' ? !uploadedFile : !carouselCards[0].visual)" class="bg-[#ed3875] text-white px-8 py-3 rounded-lg font-bold hover:opacity-90 transition-opacity disabled:opacity-50 disabled:cursor-not-allowed shadow-md">
            Check Compliance &gt;
          </button>
        </div>
      </div>
    </div>

    <!-- STEP 3: Processing State -->
    <div v-if="currentStep === 3" class="border border-border rounded-xl bg-card p-12 flex flex-col items-center justify-center text-center animate-in fade-in zoom-in-95">
       <IconLoader2 class="animate-spin text-primary mb-6" :size="48" />
       <h2 class="text-xl font-bold text-primary mb-2">AI is analyzing your ad...</h2>
       <p class="text-muted-foreground">Checking Meta policies for health claims, sensationalism, and visual compliance.</p>
       
       <div class="mt-8 space-y-3 w-full max-w-md">
         <div class="bg-primary/5 border border-primary/20 rounded-lg p-3 flex items-center justify-between">
           <span class="text-sm font-medium">Analyzing wording & copy</span>
           <IconCheck class="text-primary" :size="16" />
         </div>
         <div class="bg-muted/50 border border-border rounded-lg p-3 flex items-center justify-between opacity-70">
           <span class="text-sm font-medium">Scanning visual elements</span>
           <IconLoader2 class="text-muted-foreground animate-spin" :size="16" />
         </div>
       </div>
    </div>

    <!-- STEP 4: Results Dashboard -->
    <div v-if="currentStep === 4" class="animate-in fade-in slide-in-from-bottom-4">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h2 class="text-2xl font-bold text-foreground">Compliance Report</h2>
          <p class="text-muted-foreground mt-1">Review the issues found within your ad creative and copy.</p>
        </div>
        <button @click="currentStep = 1" class="text-sm font-medium text-primary hover:underline">Start New Check</button>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
        
        <!-- Left: Context/Preview (2 cols) -->
        <div class="col-span-2 space-y-4">
          <div class="border border-border rounded-xl bg-card p-4">
            <!-- Simulated Facebook Ad Card -->
            <div class="flex items-center gap-2 mb-3">
              <div class="w-8 h-8 rounded-full bg-muted"></div>
              <div>
                <div class="text-xs font-bold text-foreground">Your Brand Name</div>
                <div class="text-[10px] text-muted-foreground">Sponsored</div>
              </div>
            </div>
            <p class="text-sm text-foreground mb-3">{{ adCopy.primaryText || 'Lose 20kg in 30 days - GUARANTEED!' }}</p>
            <div class="w-full aspect-square bg-primary/10 rounded overflow-hidden flex items-center justify-center text-primary font-bold shadow-inner">
              [Uploaded Image Preview]
            </div>
            <div class="bg-muted/30 p-3 mt-2 border border-border flex justify-between items-center">
              <div>
                <p class="text-xs text-muted-foreground uppercase">EXAMPLE.COM</p>
                <p class="text-sm font-bold">{{ adCopy.headline || '100% Natural Solution' }}</p>
                <p class="text-xs text-muted-foreground truncate w-48">{{ adCopy.description || 'No exercise or diet needed.' }}</p>
              </div>
               <button class="bg-secondary text-secondary-foreground text-xs font-bold px-3 py-1 rounded">Learn More</button>
            </div>
          </div>
        </div>

        <!-- Right: Results Checklist (3 cols) -->
        <div class="col-span-3 space-y-4">
          
          <div class="border border-border rounded-xl bg-card overflow-hidden">
             <div class="bg-destructive/10 border-b border-destructive/20 p-4 flex items-center justify-between">
                <div class="flex items-center gap-2 text-destructive font-bold">
                  <IconAlertCircle :size="20" /> 
                  Meta Compliance Issues
                </div>
                <div class="bg-destructive text-destructive-foreground text-xs font-bold px-2 py-0.5 rounded-full">Failed: High Risk</div>
             </div>
             
             <div class="p-6 space-y-6">
                <!-- Issue 1 -->
                <div>
                  <div class="flex items-start gap-2">
                    <IconAlertCircle class="text-destructive mt-0.5" :size="16" />
                    <div>
                      <h4 class="font-bold text-sm text-foreground">Misleading Health Claims (Copy)</h4>
                      <p class="text-sm text-muted-foreground mt-1">
                        The primary text uses extreme weight loss claims ("Lose 20kg in 30 days - GUARANTEED!"). Meta severely restricts promises of guaranteed physical results without robust scientific backings or disclaimers.
                      </p>
                    </div>
                  </div>
                </div>

                <!-- Issue 2 -->
                <div>
                  <div class="flex items-start gap-2">
                    <IconAlertCircle class="text-destructive mt-0.5" :size="16" />
                    <div>
                      <h4 class="font-bold text-sm text-foreground">Unrealistic Body Standards (Visual)</h4>
                      <p class="text-sm text-muted-foreground mt-1">
                        The uploaded image appears to show a "before and after" scenario. Meta prohibits before-and-after images that promote sudden weight loss or invoke negative self-perception.
                      </p>
                    </div>
                  </div>
                </div>
             </div>
          </div>
          
          <div class="bg-muted p-4 rounded-xl border border-border flex items-start gap-3">
            <IconCheck class="text-muted-foreground shrink-0 mt-0.5" />
            <p class="text-sm text-muted-foreground">No trademark infringements or prohibited profanity detected in the ad copy or visual text overlay.</p>
          </div>

        </div>
      </div>
    </div>

  </div>
</template>

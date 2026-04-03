<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
      <div>
        <el-button @click="$router.push('/domains')" text><el-icon><ArrowLeft /></el-icon>返回</el-button>
        <h2 style="display:inline;margin-left:8px">內容編輯 - {{ domainName }}</h2>
        <el-tag style="margin-left:8px">{{ kwType }}</el-tag>
      </div>
      <el-button type="primary" :loading="saving" @click="save">保存內容</el-button>
    </div>

    <el-form :model="form" label-width="120px" v-loading="loading">
      <el-card style="margin-bottom:16px">
        <template #header>SEO基礎（必填）</template>
        <el-form-item label="目標關鍵詞"><el-input v-model="form.target_keyword" /></el-form-item>
        <el-form-item label="Page Title"><el-input v-model="form.page_title" maxlength="70" show-word-limit /></el-form-item>
        <el-form-item label="Meta Desc"><el-input v-model="form.meta_desc" type="textarea" :rows="3" maxlength="160" show-word-limit /></el-form-item>
        <el-form-item label="H1"><el-input v-model="form.h1" /></el-form-item>
      </el-card>

      <el-card style="margin-bottom:16px">
        <template #header>品牌 / 首屏</template>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="品牌名稱"><el-input v-model="form.brand_name" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="品牌色"><el-color-picker v-model="form.brand_color" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="CTA按鈕"><el-input v-model="form.cta_text" placeholder="如：立即免費註冊" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="CTA副文字"><el-input v-model="form.cta_sub" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="Hero標題"><el-input v-model="form.hero_title" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="Hero副標題"><el-input v-model="form.hero_subtitle" /></el-form-item></el-col>
        </el-row>
      </el-card>

      <el-card style="margin-bottom:16px">
        <template #header>特色賣點（3條）</template>
        <el-row :gutter="12" v-for="i in [1,2,3]" :key="i" style="margin-bottom:8px">
          <el-col :span="3"><el-input v-model="form[`feature_${i}_icon`]" placeholder="圖標" /></el-col>
          <el-col :span="7"><el-input v-model="form[`feature_${i}_title`]" placeholder="標題" /></el-col>
          <el-col :span="14"><el-input v-model="form[`feature_${i}_desc`]" placeholder="描述" /></el-col>
        </el-row>
      </el-card>

      <el-card style="margin-bottom:16px">
        <template #header>正文內容</template>
        <el-form-item label="簡介"><el-input v-model="form.intro_text" type="textarea" :rows="3" placeholder="150-200字簡介" /></el-form-item>
        <el-form-item label="正文"><el-input v-model="form.body_content" type="textarea" :rows="10" placeholder="Markdown格式，500-1500字" /></el-form-item>
        <el-form-item label="結論"><el-input v-model="form.conclusion" type="textarea" :rows="3" /></el-form-item>
      </el-card>

      <el-card style="margin-bottom:16px">
        <template #header>FAQ</template>
        <el-form-item label="FAQ標題"><el-input v-model="form.faq_title" /></el-form-item>
        <el-input v-model="form.faq_items" type="textarea" :rows="6" placeholder='[{"q":"問題","a":"回答"}]' />
      </el-card>

      <el-card style="margin-bottom:16px">
        <template #header>EEAT信號</template>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="作者"><el-input v-model="form.author_name" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="頭銜"><el-input v-model="form.author_title" /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="簡介"><el-input v-model="form.author_bio" /></el-form-item></el-col>
        </el-row>
        <el-form-item label="信任徽章"><el-input v-model="form.trust_badges" placeholder='["SSL安全加密","PAGCOR合法牌照"]' /></el-form-item>
        <el-form-item label="廣告聲明"><el-input v-model="form.disclosure" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="免責聲明"><el-input v-model="form.disclaimer" type="textarea" :rows="2" /></el-form-item>
      </el-card>

      <el-card style="margin-bottom:16px">
        <template #header>
          <span>{{ kwType }} 專屬擴展欄位</span>
        </template>
        <el-input v-model="form.extra_data" type="textarea" :rows="6" :placeholder="extraPlaceholder" />
        <div style="margin-top:8px;color:#909399;font-size:12px">
          根據詞類「{{ kwType }}」填寫對應的擴展欄位，JSON格式
        </div>
      </el-card>
    </el-form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/api'
import { ElMessage } from 'element-plus'

const route = useRoute()
const domainId = route.params.id
const domainName = ref('')
const kwType = ref('')
const loading = ref(false)
const saving = ref(false)
const form = ref({
  keyword_type: '', target_keyword: '', page_title: '', meta_desc: '', h1: '',
  brand_name: '', brand_color: '#1976D2', cta_text: '', cta_sub: '',
  hero_title: '', hero_subtitle: '',
  feature_1_icon: '', feature_1_title: '', feature_1_desc: '',
  feature_2_icon: '', feature_2_title: '', feature_2_desc: '',
  feature_3_icon: '', feature_3_title: '', feature_3_desc: '',
  intro_text: '', body_content: '', conclusion: '',
  faq_title: '常見問題', faq_items: '[]', extra_data: '{}',
  author_name: '', author_title: '', author_bio: '',
  trust_badges: '[]', disclosure: '', disclaimer: '',
})

const extraPlaceholders = {
  brand: '{"overall_rating":"9.2","bonus_amount":"NT$5000","game_count":"3000+","min_deposit":"NT$100","withdrawal_speed":"5分鐘","license":"菲律賓PAGCOR"}',
  game: '{"game_name":"百家樂","game_type":"真人","rtp":"98.76%","provider":"Evolution Gaming","min_bet":"NT$10"}',
  sports: '{"sport_type":"棒球","league":"MLB","event_name":"世界大賽","odds_type":"歐洲盤"}',
  generic: '{"ranking_list":[{"name":"","rating":"","bonus":""}],"comparison_table":""}',
  promo: '{"promo_type":"首存優惠","bonus_percentage":"100%","max_bonus":"NT$5000","wagering":"30x","valid_until":""}',
  payment: '{"payment_methods":["超商","ATM","USDT"],"min_deposit":"NT$100","processing_time":"即時","fees":"免費"}',
  affiliate: '{"commission_type":"CPA","commission_rate":"40%","payment_cycle":"月結","min_payout":"NT$3000"}',
  strategy: '{"strategy_name":"馬丁格爾","game_type":"百家樂","difficulty":"中等","expected_roi":""}',
  app: '{"platform":"iOS/Android","app_size":"45MB","version":"2.1","download_url":""}',
  register: '{"registration_steps":["步驟1","步驟2","步驟3"],"required_docs":"身分證","verification_time":"5分鐘"}',
  region: '{"city":"台北","local_features":"","nearby_casinos":"","local_regulations":""}',
  credit: '{"credit_type":"信用版","credit_limit":"NT$50000","settlement":"週結","risk_level":"高"}',
  live: '{"dealer_type":"真人荷官","stream_quality":"1080p","tables_count":"200+","languages":["中文","英文"]}',
  community: '{"platform":"PTT","post_count":"","sentiment":"正面","key_opinions":""}',
  terms: '{"term":"百家樂","definition":"","related_terms":[],"difficulty":"入門"}'
}

const extraPlaceholder = computed(() => extraPlaceholders[kwType.value] || '{}')

async function loadContent() {
  loading.value = true
  try {
    const res = await api.getContent(domainId)
    const data = res.data
    domainName.value = data.domain || ''
    kwType.value = data.keyword_type || ''
    if (data.content) {
      Object.assign(form.value, data.content)
    }
  } catch (e) {
    ElMessage.error('載入失敗: ' + e.message)
  } finally { loading.value = false }
}

async function save() {
  saving.value = true
  try {
    await api.saveContent(domainId, form.value)
    ElMessage.success('內容已保存')
  } catch (e) {
    ElMessage.error(e.message)
  } finally { saving.value = false }
}

onMounted(loadContent)
</script>

<template>
  <div>
    <h2>系統設定</h2>
    <el-card v-loading="loading">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基礎配置" name="basic">
          <el-form :model="settings" label-width="180px" style="max-width:600px">
            <el-form-item label="Hugo路徑"><el-input v-model="settings.hugo_path" /></el-form-item>
            <el-form-item label="站點根目錄"><el-input v-model="settings.sites_root" /></el-form-item>
            <el-form-item label="預設服務器IP"><el-input v-model="settings.default_server_ip" /></el-form-item>
            <el-form-item label="預設SSH用戶"><el-input v-model="settings.default_server_user" /></el-form-item>
            <el-form-item label="追蹤器URL"><el-input v-model="settings.tracker_url" placeholder="https://tracker.yourdomain.com" /></el-form-item>
            <el-form-item label="CC限速(次/分鐘)"><el-input-number v-model.number="ccLimit" :min="10" :max="1000" /></el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="API密鑰" name="api">
          <el-form :model="settings" label-width="180px" style="max-width:600px">
            <el-form-item label="Claude API Key">
              <el-input v-model="settings.claude_api_key" type="password" show-password placeholder="sk-ant-..." />
              <div style="color:#909399;font-size:12px">用於AI內容生成</div>
            </el-form-item>
            <el-form-item label="Ahrefs API Key"><el-input v-model="settings.ahrefs_api_key" type="password" show-password /></el-form-item>
            <el-form-item label="SerpAPI Key"><el-input v-model="settings.serpapi_key" type="password" show-password /></el-form-item>
            <el-form-item label="Bing WMT Key"><el-input v-model="settings.bing_wmt_key" type="password" show-password /></el-form-item>
            <el-divider />
            <el-form-item label="修改密碼">
              <el-row :gutter="12">
                <el-col :span="8"><el-input v-model="oldPw" type="password" placeholder="舊密碼" /></el-col>
                <el-col :span="8"><el-input v-model="newPw" type="password" placeholder="新密碼(至少6位)" /></el-col>
                <el-col :span="4"><el-button @click="changePw">修改</el-button></el-col>
              </el-row>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <el-tab-pane label="功能開關" name="features">
          <el-form :model="settings" label-width="220px" style="max-width:600px">
            <el-form-item label="構建後自動推送IndexNow">
              <el-switch v-model="settings.indexnow_auto" active-value="1" inactive-value="0" />
            </el-form-item>
            <el-form-item label="內容自動更新調度">
              <el-switch v-model="settings.content_refresh_on" active-value="1" inactive-value="0" />
            </el-form-item>
            <el-form-item label="城市詞矩陣">
              <el-switch v-model="settings.city_matrix_on" active-value="1" inactive-value="0" />
            </el-form-item>
            <el-form-item label="OG圖片動態生成">
              <el-switch v-model="settings.og_image_enabled" active-value="1" inactive-value="0" />
            </el-form-item>
            <el-form-item label="支撐頁面自動生成">
              <el-switch v-model="settings.support_pages_auto" active-value="1" inactive-value="0" />
            </el-form-item>
            <el-form-item label="CSS噪聲注入(反指紋)">
              <el-switch v-model="settings.css_noise_enabled" active-value="1" inactive-value="0" />
            </el-form-item>
            <el-form-item label="Schema欄位順序隨機化">
              <el-switch v-model="settings.schema_shuffle" active-value="1" inactive-value="0" />
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>

      <div style="margin-top:24px">
        <el-button type="primary" @click="save" :loading="saving">保存設定</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const activeTab = ref('basic')
const settings = ref({})
const loading = ref(false)
const saving = ref(false)
const oldPw = ref('')
const newPw = ref('')

const ccLimit = computed({
  get: () => parseInt(settings.value.cc_limit_per_min) || 120,
  set: (v) => { settings.value.cc_limit_per_min = String(v) }
})

async function load() {
  loading.value = true
  try {
    const res = await api.getSettings()
    settings.value = res.data || {}
  } finally { loading.value = false }
}

async function save() {
  saving.value = true
  try {
    await api.saveSettings(settings.value)
    ElMessage.success('設定已保存')
  } catch (e) { ElMessage.error(e.message) } finally { saving.value = false }
}

async function changePw() {
  if (!newPw.value || newPw.value.length < 6) { ElMessage.warning('新密碼至少6位'); return }
  try {
    const res = await api.changePassword({ old_password: oldPw.value, new_password: newPw.value })
    localStorage.setItem('token', res.data.token)
    ElMessage.success('密碼已修改')
    oldPw.value = ''; newPw.value = ''
  } catch (e) { ElMessage.error(e.message) }
}

onMounted(load)
</script>

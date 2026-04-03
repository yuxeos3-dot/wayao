<template>
  <div>
    <h2>排名追蹤</h2>
    <el-alert type="info" :closable="false" style="margin-bottom:16px">
      排名數據由 Python 腳本定期採集，此頁面顯示歷史記錄。使用 <code>scripts/check_rankings.py</code> 執行排名檢查。
    </el-alert>

    <el-row :gutter="12" style="margin-bottom:12px">
      <el-col :span="6">
        <el-select v-model="domainId" placeholder="選擇域名" filterable clearable @change="load">
          <el-option v-for="d in domains" :key="d.id" :label="d.domain" :value="d.id" />
        </el-select>
      </el-col>
    </el-row>

    <el-table :data="rankings" v-loading="loading" size="small">
      <el-table-column prop="keyword" label="關鍵詞" min-width="200" />
      <el-table-column prop="position" label="排名" width="80">
        <template #default="{row}">
          <el-tag :type="row.position<=3?'success':row.position<=10?'warning':'info'" size="small">{{ row.position || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="url" label="URL" min-width="200" />
      <el-table-column prop="checked_at" label="檢查時間" width="170" />
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const domains = ref([])
const rankings = ref([])
const loading = ref(false)
const domainId = ref('')

async function load() {
  if (!domainId.value) { rankings.value = []; return }
  loading.value = true
  try {
    // ranking data is collected by Python script, query from ranking_history
    const res = await api.get(`/stats/clicks`, { site_id: domains.value.find(d => d.id === domainId.value)?.domain })
    rankings.value = res.data || []
  } catch (e) {
    console.warn('Rankings load failed:', e)
  } finally { loading.value = false }
}

onMounted(async () => {
  try { const res = await api.listDomains({ size: 500 }); domains.value = res.data || [] } catch {}
})
</script>

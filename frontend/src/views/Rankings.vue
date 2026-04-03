<template>
  <div>
<<<<<<< HEAD
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
      <el-table-column prop="rank" label="排名" width="80">
        <template #default="{row}">
          <el-tag :type="row.rank<=3?'success':row.rank<=10?'warning':'info'" size="small">{{ row.rank || '-' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="engine" label="引擎" width="80" />
      <el-table-column prop="checked_at" label="檢查時間" width="170" />
    </el-table>
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">排名追踪</h1>
        <p class="page-subtitle">关键词排名历史数据</p>
      </div>
    </div>

    <el-alert type="info" :closable="false" style="margin-bottom:20px;border-radius:var(--wayao-radius-sm)">
      排名数据由 Python 脚本定期采集, 此页面展示历史记录。执行 <code>scripts/check_rankings.py</code> 进行排名检查。
    </el-alert>

    <div class="filter-bar">
      <el-select v-model="domainId" placeholder="选择域名" filterable clearable @change="load" style="width:280px">
        <el-option v-for="d in domains" :key="d.id" :label="d.domain" :value="d.id" />
      </el-select>
    </div>

    <div class="section-card">
      <div class="section-body" :style="rankings.length ? {padding:0} : {}">
        <el-table :data="rankings" v-loading="loading" size="small" v-if="rankings.length">
          <el-table-column prop="keyword" label="关键词" min-width="200">
            <template #default="{row}"><span style="font-weight:500">{{ row.keyword }}</span></template>
          </el-table-column>
          <el-table-column prop="rank" label="排名" width="80">
            <template #default="{row}">
              <el-tag :type="row.rank<=3?'success':row.rank<=10?'warning':'info'" size="small">{{ row.rank || '-' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="engine" label="引擎" width="80" />
          <el-table-column prop="checked_at" label="检查时间" width="170" />
        </el-table>
        <el-empty v-else :description="domainId ? '暂无排名数据' : '请选择域名查看'" :image-size="60" />
      </div>
    </div>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
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
<<<<<<< HEAD
  try {
    // ranking data is collected by Python script, query from ranking_history
    const res = await api.get('/rankings', { domain_id: domainId.value })
    rankings.value = res.data || []
  } catch (e) {
    console.warn('Rankings load failed:', e)
  } finally { loading.value = false }
=======
  try { const res = await api.get('/rankings', { domain_id: domainId.value }); rankings.value = res.data || [] }
  catch {} finally { loading.value = false }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
}

onMounted(async () => {
  try { const res = await api.listDomains({ size: 500 }); domains.value = res.data || [] } catch {}
})
</script>

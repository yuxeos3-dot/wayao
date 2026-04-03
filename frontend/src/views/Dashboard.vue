<template>
  <div>
    <h2>數據看板</h2>
    <el-row :gutter="20" style="margin-bottom:20px">
      <el-col :span="6" v-for="card in cards" :key="card.label">
        <el-card shadow="hover">
          <div style="font-size:14px;color:#909399">{{ card.label }}</div>
          <div style="font-size:28px;font-weight:bold;margin:8px 0">{{ card.value }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>今日點擊分佈</template>
          <div v-if="Object.keys(stats.clicks_by_action||{}).length">
            <div v-for="(cnt, action) in stats.clicks_by_action" :key="action" style="display:flex;justify-content:space-between;padding:4px 0">
              <span>{{ action }}</span><el-tag>{{ cnt }}</el-tag>
            </div>
          </div>
          <el-empty v-else description="暫無數據" :image-size="60" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>Top 域名 (7天)</template>
          <el-table :data="stats.top_domains||[]" size="small">
            <el-table-column prop="domain" label="域名" />
            <el-table-column prop="clicks" label="點擊" width="80" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top:20px">
      <template #header>健康告警</template>
      <el-table :data="alerts" size="small" v-if="alerts.length">
        <el-table-column prop="domain" label="域名" />
        <el-table-column prop="status" label="狀態">
          <template #default="{row}"><el-tag type="danger">{{ row.status }}</el-tag></template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="一切正常" :image-size="60" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api'

const stats = ref({})
const alerts = ref([])
const cards = computed(() => [
  { label: '總域名', value: stats.value.total_domains || 0 },
  { label: '已上線', value: stats.value.active_domains || 0 },
  { label: '今日點擊', value: stats.value.today_clicks || 0 },
  { label: '今日PV', value: stats.value.today_pv || 0 },
])

onMounted(async () => {
  try {
    const res = await api.overview()
    stats.value = res.data || {}
  } catch {}
  try {
    const res = await api.get('/health-check/alerts')
    alerts.value = res.data || []
  } catch {}
})
</script>

<template>
  <div>
<<<<<<< HEAD
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
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">数据看板</h1>
        <p class="page-subtitle">站群运营数据概览</p>
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stat-grid stat-grid-4">
      <div class="stat-card" v-for="card in cards" :key="card.label">
        <div class="stat-label">{{ card.label }}</div>
        <div class="stat-value">{{ card.value }}</div>
      </div>
    </div>

    <!-- Charts Row -->
    <el-row :gutter="20" style="margin-bottom:24px">
      <el-col :span="12">
        <div class="section-card">
          <div class="section-header">今日点击分布</div>
          <div class="section-body">
            <div v-if="Object.keys(stats.clicks_by_action||{}).length">
              <div v-for="(cnt, action) in stats.clicks_by_action" :key="action" class="action-row">
                <span class="action-name">{{ action }}</span>
                <div class="action-bar-wrap">
                  <div class="action-bar" :style="{width: barWidth(cnt) + '%'}"></div>
                </div>
                <span class="action-count">{{ cnt }}</span>
              </div>
            </div>
            <el-empty v-else description="暂无数据" :image-size="60" />
          </div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="section-card">
          <div class="section-header">热门站点 (近7天)</div>
          <div class="section-body" style="padding:0">
            <el-table :data="stats.top_domains||[]" size="small" :show-header="false" style="border-radius:0">
              <el-table-column prop="domain" />
              <el-table-column prop="clicks" width="80" align="right">
                <template #default="{row}"><span style="font-weight:600;color:var(--wayao-accent)">{{ row.clicks }}</span></template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- Health Alerts -->
    <div class="section-card">
      <div class="section-header">
        <span>健康告警</span>
        <el-tag v-if="alerts.length" type="danger" size="small">{{ alerts.length }}</el-tag>
      </div>
      <div class="section-body" :style="alerts.length ? {padding:0} : {}">
        <el-table :data="alerts" size="small" v-if="alerts.length">
          <el-table-column prop="domain" label="域名" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{row}"><el-tag type="danger" size="small">{{ row.status }}</el-tag></template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="一切正常" :image-size="60" />
      </div>
    </div>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import api from '@/api'

const stats = ref({})
const alerts = ref([])
const cards = computed(() => [
<<<<<<< HEAD
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
=======
  { label: '总域名数', value: stats.value.total_domains || 0 },
  { label: '已上线', value: stats.value.active_domains || 0 },
  { label: '今日点击', value: stats.value.today_clicks || 0 },
  { label: '今日PV', value: stats.value.today_pv || 0 },
])

const maxAction = computed(() => {
  const actions = stats.value.clicks_by_action || {}
  return Math.max(1, ...Object.values(actions))
})

function barWidth(cnt) {
  return Math.round((cnt / maxAction.value) * 100)
}

onMounted(async () => {
  try { const res = await api.overview(); stats.value = res.data || {} } catch {}
  try { const res = await api.get('/health-check/alerts'); alerts.value = res.data || [] } catch {}
})
</script>

<style scoped>
.action-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
}
.action-row + .action-row { border-top: 1px solid var(--wayao-border); }
.action-name { font-size: 14px; font-weight: 500; width: 80px; flex-shrink: 0; }
.action-bar-wrap { flex: 1; height: 8px; background: rgba(0,0,0,0.04); border-radius: 4px; overflow: hidden; }
.action-bar { height: 100%; background: var(--wayao-accent); border-radius: 4px; transition: width 0.6s ease; }
.action-count { font-size: 14px; font-weight: 600; color: var(--wayao-text); width: 48px; text-align: right; }
</style>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31

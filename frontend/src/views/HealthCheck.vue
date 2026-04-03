<template>
  <div>
<<<<<<< HEAD
    <h2>健康檢查</h2>
    <el-card style="margin-bottom:16px">
      <template #header>告警列表</template>
      <el-table :data="alerts" v-loading="loading" size="small">
        <el-table-column prop="domain" label="域名" />
        <el-table-column prop="status" label="狀態">
          <template #default="{row}"><el-tag type="danger">{{ row.status }}</el-tag></template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!alerts.length && !loading" description="一切正常" />
    </el-card>
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">健康检查</h1>
        <p class="page-subtitle">监控站点运行状态, 自动检测异常</p>
      </div>
    </div>

    <div class="section-card">
      <div class="section-header">
        <span>告警列表</span>
        <el-tag v-if="alerts.length" type="danger" size="small">{{ alerts.length }} 个异常</el-tag>
      </div>
      <div class="section-body" :style="alerts.length ? {padding:0} : {}">
        <el-table :data="alerts" v-loading="loading" size="small" v-if="alerts.length">
          <el-table-column prop="domain" label="域名">
            <template #default="{row}"><span style="font-weight:500">{{ row.domain }}</span></template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{row}"><el-tag type="danger" size="small">{{ row.status === 'error' ? '异常' : '下线' }}</el-tag></template>
          </el-table-column>
        </el-table>
        <el-empty v-else description="一切正常, 所有站点运行良好" :image-size="80" />
      </div>
    </div>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const alerts = ref([])
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try { const res = await api.getHealthAlerts(); alerts.value = res.data || [] } catch {} finally { loading.value = false }
})
</script>

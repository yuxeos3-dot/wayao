<template>
  <div>
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

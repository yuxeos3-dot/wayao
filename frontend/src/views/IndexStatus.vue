<template>
  <div>
    <h2>收錄監控</h2>
    <el-row :gutter="12" style="margin-bottom:12px">
      <el-col :span="6">
        <el-select v-model="domainId" placeholder="選擇域名" filterable @change="checkIndex">
          <el-option v-for="d in domains" :key="d.id" :label="d.domain" :value="d.id" />
        </el-select>
      </el-col>
      <el-col :span="4"><el-button type="primary" @click="batchCheck" :loading="loading">批量檢查</el-button></el-col>
    </el-row>

    <el-card v-if="indexResult">
      <template #header>收錄狀態</template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="域名">{{ indexResult.domain }}</el-descriptions-item>
        <el-descriptions-item label="狀態">{{ indexResult.status }}</el-descriptions-item>
      </el-descriptions>
      <el-alert type="info" style="margin-top:12px" :closable="false">
        {{ indexResult.tip }}
      </el-alert>
    </el-card>

    <el-card style="margin-top:16px">
      <template #header>IndexNow 推送記錄</template>
      <el-table :data="submissions" size="small" v-if="domainId">
        <el-table-column prop="url" label="URL" min-width="250" />
        <el-table-column prop="engine" label="引擎" width="100" />
        <el-table-column prop="status" label="狀態" width="80">
          <template #default="{row}"><el-tag :type="row.status==='submitted'?'success':'info'" size="small">{{ row.status }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" label="時間" width="170" />
      </el-table>
      <el-empty v-else description="選擇域名查看推送記錄" />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const domains = ref([])
const domainId = ref('')
const indexResult = ref(null)
const submissions = ref([])
const loading = ref(false)

async function checkIndex() {
  if (!domainId.value) return
  try {
    const res = await api.get(`/index-status/${domainId.value}`)
    indexResult.value = res.data
    const sub = await api.getIndexNowRecords(domainId.value)
    submissions.value = sub.data || []
  } catch (e) { ElMessage.error(e.message) }
}

async function batchCheck() {
  loading.value = true
  try {
    await api.post('/index-status/batch')
    ElMessage.success('批量檢查已啟動')
  } catch (e) { ElMessage.error(e.message) }
  finally { loading.value = false }
}

onMounted(async () => {
  try { const res = await api.listDomains({ size: 500 }); domains.value = res.data || [] } catch {}
})
</script>

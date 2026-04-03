<template>
  <div>
    <h2>構建部署</h2>
    <el-row :gutter="12" style="margin-bottom:16px">
      <el-col :span="4"><el-select v-model="filter.status" placeholder="狀態" clearable @change="load"><el-option v-for="s in ['draft','built','active','error']" :key="s" :label="s" :value="s" /></el-select></el-col>
    </el-row>

    <el-table :data="list" v-loading="loading" @selection-change="s=>selected=s">
      <el-table-column type="selection" width="40" />
      <el-table-column prop="domain" label="域名" min-width="160" />
      <el-table-column prop="keyword_type" label="詞類" width="80" />
      <el-table-column prop="template_name" label="模版" width="100" />
      <el-table-column prop="has_content" label="內容" width="60">
        <template #default="{row}"><el-tag :type="row.has_content?'success':'danger'" size="small">{{ row.has_content?'有':'無' }}</el-tag></template>
      </el-table-column>
      <el-table-column prop="status" label="狀態" width="80">
        <template #default="{row}"><el-tag :type="statusColor(row.status)" size="small">{{ row.status }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="280">
        <template #default="{row}">
          <el-button size="small" type="primary" @click="buildOne(row.id)" :loading="building[row.id]">構建</el-button>
          <el-button size="small" type="success" @click="deployOne(row.id)" :loading="deploying[row.id]">部署</el-button>
          <el-button size="small" @click="showLog(row.id)">日誌</el-button>
          <el-button size="small" @click="api.exportDomain(row.id)">匯出</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="margin-top:16px">
      <el-button type="primary" :disabled="!selected.length" @click="batchBuild">批量構建 ({{ selected.length }})</el-button>
      <el-button type="success" :disabled="!selected.length" @click="batchDeploy">批量部署 ({{ selected.length }})</el-button>
    </div>

    <el-dialog v-model="logVisible" title="構建日誌" width="700">
      <el-table :data="logs" size="small">
        <el-table-column prop="status" label="狀態" width="80">
          <template #default="{row}"><el-tag :type="row.status==='success'?'success':'danger'" size="small">{{ row.status }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="duration_ms" label="耗時(ms)" width="100" />
        <el-table-column prop="created_at" label="時間" width="170" />
        <el-table-column prop="output" label="輸出" min-width="200">
          <template #default="{row}"><pre style="max-height:100px;overflow:auto;font-size:11px;margin:0">{{ row.log_output }}</pre></template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(false)
const selected = ref([])
const filter = ref({ status: '' })
const building = reactive({})
const deploying = reactive({})
const logVisible = ref(false)
const logs = ref([])

function statusColor(s) {
  return { active: 'success', built: 'warning', draft: 'info', error: 'danger' }[s] || 'info'
}

async function load() {
  loading.value = true
  try {
    const res = await api.listDomains({ ...filter.value, size: 200 })
    list.value = res.data || []
  } finally { loading.value = false }
}

async function buildOne(id) {
  building[id] = true
  try { await api.buildSite(id); ElMessage.success('構建完成'); load() } catch (e) { ElMessage.error(e.message) } finally { building[id] = false }
}

async function deployOne(id) {
  deploying[id] = true
  try { await api.deploySite(id); ElMessage.success('部署完成'); load() } catch (e) { ElMessage.error(e.message) } finally { deploying[id] = false }
}

async function batchBuild() {
  try { await api.batchDomainOp({ ids: selected.value.map(r => r.id), action: 'build' }); ElMessage.success('批量構建完成'); load() } catch (e) { ElMessage.error(e.message) }
}

async function batchDeploy() {
  try { await api.batchDomainOp({ ids: selected.value.map(r => r.id), action: 'deploy' }); ElMessage.success('批量部署完成'); load() } catch (e) { ElMessage.error(e.message) }
}

async function showLog(id) {
  try { const res = await api.getBuildLog(id); logs.value = res.data || []; logVisible.value = true } catch (e) { ElMessage.error(e.message) }
}

onMounted(load)
</script>

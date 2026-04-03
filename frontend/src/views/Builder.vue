<template>
  <div>
<<<<<<< HEAD
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
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">构建部署</h1>
        <p class="page-subtitle">Hugo 站点构建与远程部署管理</p>
      </div>
      <div style="display:flex;gap:12px">
        <el-button :disabled="!selected.length" @click="batchBuild">批量构建 ({{ selected.length }})</el-button>
        <el-button type="success" :disabled="!selected.length" @click="batchDeploy">批量部署 ({{ selected.length }})</el-button>
      </div>
    </div>

    <div class="filter-bar">
      <el-select v-model="filter.status" placeholder="状态筛选" clearable @change="load" style="width:140px">
        <el-option v-for="s in [{v:'draft',l:'草稿'},{v:'built',l:'已构建'},{v:'active',l:'已上线'},{v:'error',l:'异常'}]" :key="s.v" :label="s.l" :value="s.v" />
      </el-select>
    </div>

    <div class="section-card">
      <div class="section-body" style="padding:0">
        <el-table :data="list" v-loading="loading" @selection-change="s=>selected=s">
          <el-table-column type="selection" width="40" />
          <el-table-column prop="domain" label="域名" min-width="180">
            <template #default="{row}"><span style="font-weight:500">{{ row.domain }}</span></template>
          </el-table-column>
          <el-table-column prop="keyword_type" label="词类" width="90">
            <template #default="{row}"><el-tag size="small">{{ row.keyword_type }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="template_name" label="模版" width="100" />
          <el-table-column prop="has_content" label="内容" width="60">
            <template #default="{row}"><el-tag :type="row.has_content?'success':'danger'" size="small">{{ row.has_content?'有':'无' }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{row}"><el-tag :type="statusColor(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag></template>
          </el-table-column>
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{row}">
              <el-button size="small" type="primary" @click="buildOne(row.id)" :loading="building[row.id]">构建</el-button>
              <el-button size="small" type="success" @click="deployOne(row.id)" :loading="deploying[row.id]">部署</el-button>
              <el-button size="small" @click="showLog(row.id)">日志</el-button>
              <el-button size="small" @click="api.exportDomain(row.id)">导出</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-dialog v-model="logVisible" title="构建日志" width="720">
      <el-table :data="logs" size="small">
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{row}"><el-tag :type="row.status==='success'?'success':'danger'" size="small">{{ row.status==='success'?'成功':'失败' }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="duration_ms" label="耗时(ms)" width="100" />
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column prop="output" label="输出" min-width="200">
          <template #default="{row}"><pre style="max-height:100px;overflow:auto;font-size:11px;margin:0;white-space:pre-wrap">{{ row.log_output }}</pre></template>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
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
<<<<<<< HEAD

async function load() {
  loading.value = true
  try {
    const res = await api.listDomains({ ...filter.value, size: 200 })
    list.value = res.data || []
  } finally { loading.value = false }
=======
function statusLabel(s) {
  return { active: '已上线', built: '已构建', draft: '草稿', error: '异常' }[s] || s
}

async function load() {
  loading.value = true
  try { const res = await api.listDomains({ ...filter.value, size: 200 }); list.value = res.data || [] } finally { loading.value = false }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
}

async function buildOne(id) {
  building[id] = true
<<<<<<< HEAD
  try { await api.buildSite(id); ElMessage.success('構建完成'); load() } catch (e) { ElMessage.error(e.message) } finally { building[id] = false }
=======
  try { await api.buildSite(id); ElMessage.success('构建完成'); load() } catch (e) { ElMessage.error(e.message) } finally { building[id] = false }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
}

async function deployOne(id) {
  deploying[id] = true
  try { await api.deploySite(id); ElMessage.success('部署完成'); load() } catch (e) { ElMessage.error(e.message) } finally { deploying[id] = false }
}

async function batchBuild() {
<<<<<<< HEAD
  try { await api.batchBuild({ ids: selected.value.map(r => r.id), action: 'build' }); ElMessage.success('批量構建已啟動'); setTimeout(load, 5000) } catch (e) { ElMessage.error(e.message) }
}

async function batchDeploy() {
  try { await api.batchBuild({ ids: selected.value.map(r => r.id), action: 'deploy' }); ElMessage.success('批量部署已啟動'); setTimeout(load, 5000) } catch (e) { ElMessage.error(e.message) }
=======
  try { await api.batchBuild({ ids: selected.value.map(r => r.id), action: 'build' }); ElMessage.success('批量构建已启动'); setTimeout(load, 5000) } catch (e) { ElMessage.error(e.message) }
}

async function batchDeploy() {
  try { await api.batchBuild({ ids: selected.value.map(r => r.id), action: 'deploy' }); ElMessage.success('批量部署已启动'); setTimeout(load, 5000) } catch (e) { ElMessage.error(e.message) }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
}

async function showLog(id) {
  try { const res = await api.getBuildLog(id); logs.value = res.data || []; logVisible.value = true } catch (e) { ElMessage.error(e.message) }
}

onMounted(load)
</script>

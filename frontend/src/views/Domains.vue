<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
      <h2 style="margin:0">域名管理</h2>
      <div>
        <el-button @click="batchOp('build')" :disabled="!selected.length">批量構建</el-button>
        <el-button type="primary" @click="showAdd=true">新增域名</el-button>
      </div>
    </div>

    <el-row :gutter="12" style="margin-bottom:12px">
      <el-col :span="4"><el-select v-model="filter.market" placeholder="市場" clearable @change="load"><el-option v-for="m in markets" :key="m" :label="m" :value="m" /></el-select></el-col>
      <el-col :span="4"><el-select v-model="filter.status" placeholder="狀態" clearable @change="load"><el-option v-for="s in statuses" :key="s" :label="s" :value="s" /></el-select></el-col>
      <el-col :span="4"><el-select v-model="filter.keyword_type" placeholder="詞類" clearable @change="load"><el-option v-for="t in kwTypes" :key="t" :label="t" :value="t" /></el-select></el-col>
    </el-row>

    <el-table :data="list" v-loading="loading" @selection-change="onSelect">
      <el-table-column type="selection" width="40" />
      <el-table-column prop="id" label="ID" width="55" />
      <el-table-column prop="domain" label="域名" min-width="160" />
      <el-table-column prop="keyword_type" label="詞類" width="80" />
      <el-table-column prop="primary_keyword" label="主關鍵詞" width="140" />
      <el-table-column prop="template_name" label="模版" width="100" />
      <el-table-column prop="status" label="狀態" width="80">
        <template #default="{row}">
          <el-tag :type="statusColor(row.status)" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="has_content" label="內容" width="60">
        <template #default="{row}"><el-tag :type="row.has_content?'success':'info'" size="small">{{ row.has_content?'有':'無' }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="260">
        <template #default="{row}">
          <el-button size="small" @click="$router.push(`/domains/${row.id}/content`)">內容</el-button>
          <el-button size="small" type="success" @click="build(row.id)">構建</el-button>
          <el-button size="small" @click="editRow(row)">編輯</el-button>
          <el-popconfirm title="確定刪除?" @confirm="del(row.id)">
            <template #reference><el-button size="small" type="danger">刪除</el-button></template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination style="margin-top:16px;justify-content:flex-end" v-model:current-page="page" :page-size="50" :total="total" @current-change="load" layout="total, prev, pager, next" />

    <el-dialog v-model="showAdd" :title="form.id?'編輯域名':'新增域名'" width="550">
      <el-form :model="form" label-width="100px">
        <el-form-item label="域名"><el-input v-model="form.domain" /></el-form-item>
        <el-form-item label="市場"><el-select v-model="form.market"><el-option v-for="m in markets" :key="m" :label="m" :value="m" /></el-select></el-form-item>
        <el-form-item label="詞類"><el-select v-model="form.keyword_type"><el-option v-for="t in kwTypes" :key="t" :label="t" :value="t" /></el-select></el-form-item>
        <el-form-item label="主關鍵詞"><el-input v-model="form.primary_keyword" /></el-form-item>
        <el-form-item label="跳轉URL"><el-input v-model="form.redirect_url" /></el-form-item>
        <el-form-item label="模版" v-if="form.id"><el-select v-model="form.template_id"><el-option v-for="t in templates" :key="t.id" :label="t.name" :value="t.id" /></el-select></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAdd=false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(false)
const showAdd = ref(false)
const page = ref(1)
const total = ref(0)
const selected = ref([])
const templates = ref([])
const filter = ref({ market: '', status: '', keyword_type: '' })
const form = ref({ domain: '', market: 'zh-TW', keyword_type: 'brand', primary_keyword: '', redirect_url: '' })

const markets = ['zh-TW', 'vi', 'th', 'pt-BR']
const statuses = ['draft', 'built', 'active', 'error', 'down']
const kwTypes = ['brand','game','sports','generic','promo','payment','affiliate','strategy','app','register','region','credit','live','community','terms']

function statusColor(s) {
  return { active: 'success', built: 'warning', draft: 'info', error: 'danger', down: 'danger' }[s] || 'info'
}

async function load() {
  loading.value = true
  try {
    const res = await api.listDomains({ ...filter.value, page: page.value, size: 50 })
    list.value = res.data || []
    total.value = res.total || 0
  } finally { loading.value = false }
}

function onSelect(rows) { selected.value = rows }

function editRow(row) { form.value = { ...row }; showAdd.value = true }

async function save() {
  try {
    if (form.value.id) await api.updateDomain(form.value.id, form.value)
    else await api.createDomain(form.value)
    ElMessage.success('已保存')
    showAdd.value = false
    form.value = { domain: '', market: 'zh-TW', keyword_type: 'brand', primary_keyword: '', redirect_url: '' }
    load()
  } catch (e) { ElMessage.error(e.message) }
}

async function del(id) {
  try { await api.deleteDomain(id); ElMessage.success('已刪除'); load() } catch (e) { ElMessage.error(e.message) }
}

async function build(id) {
  try { await api.buildSite(id); ElMessage.success('構建完成'); load() } catch (e) { ElMessage.error(e.message) }
}

async function batchOp(action) {
  try {
    await api.batchDomainOp({ ids: selected.value.map(r => r.id), action })
    ElMessage.success('批量操作完成')
    load()
  } catch (e) { ElMessage.error(e.message) }
}

onMounted(async () => {
  load()
  try { const res = await api.listTemplates(); templates.value = res.data || [] } catch {}
})
</script>

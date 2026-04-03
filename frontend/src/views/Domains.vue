<template>
  <div>
<<<<<<< HEAD
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
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">域名管理</h1>
        <p class="page-subtitle">管理所有站点域名及其配置</p>
      </div>
      <div style="display:flex;gap:12px">
        <el-button @click="batchOp('build')" :disabled="!selected.length">批量构建 ({{ selected.length }})</el-button>
        <el-button type="primary" @click="openAdd">新增域名</el-button>
      </div>
    </div>

    <!-- Filters -->
    <div class="filter-bar">
      <el-select v-model="filter.market" placeholder="市场" clearable @change="load" style="width:140px">
        <el-option v-for="m in markets" :key="m" :label="m" :value="m" />
      </el-select>
      <el-select v-model="filter.status" placeholder="状态" clearable @change="load" style="width:140px">
        <el-option v-for="s in statuses" :key="s.value" :label="s.label" :value="s.value" />
      </el-select>
      <el-select v-model="filter.keyword_type" placeholder="词类" clearable @change="load" style="width:140px">
        <el-option v-for="t in kwTypes" :key="t" :label="t" :value="t" />
      </el-select>
    </div>

    <!-- Table -->
    <div class="section-card">
      <div class="section-body" style="padding:0">
        <el-table :data="list" v-loading="loading" @selection-change="onSelect">
          <el-table-column type="selection" width="40" />
          <el-table-column prop="id" label="ID" width="55" />
          <el-table-column prop="domain" label="域名" min-width="180">
            <template #default="{row}"><span style="font-weight:500">{{ row.domain }}</span></template>
          </el-table-column>
          <el-table-column prop="keyword_type" label="词类" width="90">
            <template #default="{row}"><el-tag size="small">{{ row.keyword_type }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="primary_keyword" label="主关键词" width="140" />
          <el-table-column prop="template_name" label="模版" width="100" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{row}">
              <el-tag :type="statusColor(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="has_content" label="内容" width="60">
            <template #default="{row}"><el-tag :type="row.has_content?'success':'info'" size="small">{{ row.has_content?'有':'无' }}</el-tag></template>
          </el-table-column>
          <el-table-column label="操作" width="260" fixed="right">
            <template #default="{row}">
              <el-button size="small" @click="$router.push(`/sites/${row.id}/content`)">内容</el-button>
              <el-button size="small" type="success" @click="build(row.id)">构建</el-button>
              <el-button size="small" @click="editRow(row)">编辑</el-button>
              <el-popconfirm title="确定删除该域名?" @confirm="del(row.id)">
                <template #reference><el-button size="small" type="danger">删除</el-button></template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-pagination
      style="margin-top:16px;justify-content:flex-end"
      v-model:current-page="page" :page-size="50" :total="total"
      @current-change="load" layout="total, prev, pager, next"
    />

    <!-- Add/Edit Dialog -->
    <el-dialog v-model="showAdd" :title="form.id?'编辑域名':'新增域名'" width="550">
      <el-form :model="form" label-width="100px">
        <el-form-item label="域名"><el-input v-model="form.domain" placeholder="example.com" /></el-form-item>
        <el-form-item label="市场">
          <el-select v-model="form.market" style="width:100%">
            <el-option v-for="m in markets" :key="m" :label="m" :value="m" />
          </el-select>
        </el-form-item>
        <el-form-item label="词类">
          <el-select v-model="form.keyword_type" style="width:100%">
            <el-option v-for="t in kwTypes" :key="t" :label="t" :value="t" />
          </el-select>
        </el-form-item>
        <el-form-item label="主关键词"><el-input v-model="form.primary_keyword" /></el-form-item>
        <el-form-item label="跳转URL"><el-input v-model="form.redirect_url" placeholder="https://..." /></el-form-item>
        <el-form-item label="模版" v-if="form.id">
          <el-select v-model="form.template_id" style="width:100%">
            <el-option v-for="t in templates" :key="t.id" :label="t.name" :value="t.id" />
          </el-select>
        </el-form-item>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
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
<<<<<<< HEAD
const statuses = ['draft', 'built', 'active', 'error', 'down']
=======
const statuses = [
  { value: 'draft', label: '草稿' },
  { value: 'built', label: '已构建' },
  { value: 'active', label: '已上线' },
  { value: 'error', label: '异常' },
  { value: 'down', label: '下线' },
]
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
const kwTypes = ['brand','game','sports','generic','promo','payment','affiliate','strategy','app','register','region','credit','live','community','terms']

function statusColor(s) {
  return { active: 'success', built: 'warning', draft: 'info', error: 'danger', down: 'danger' }[s] || 'info'
}
<<<<<<< HEAD
=======
function statusLabel(s) {
  return { active: '已上线', built: '已构建', draft: '草稿', error: '异常', down: '下线' }[s] || s
}
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31

async function load() {
  loading.value = true
  try {
    const res = await api.listDomains({ ...filter.value, page: page.value, size: 50 })
    list.value = res.data || []
    total.value = res.total || 0
  } finally { loading.value = false }
}

function onSelect(rows) { selected.value = rows }

<<<<<<< HEAD
=======
function openAdd() {
  form.value = { domain: '', market: 'zh-TW', keyword_type: 'brand', primary_keyword: '', redirect_url: '' }
  showAdd.value = true
}

>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
function editRow(row) { form.value = { ...row }; showAdd.value = true }

async function save() {
  try {
    if (form.value.id) await api.updateDomain(form.value.id, form.value)
    else await api.createDomain(form.value)
    ElMessage.success('已保存')
    showAdd.value = false
<<<<<<< HEAD
    form.value = { domain: '', market: 'zh-TW', keyword_type: 'brand', primary_keyword: '', redirect_url: '' }
=======
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
    load()
  } catch (e) { ElMessage.error(e.message) }
}

async function del(id) {
<<<<<<< HEAD
  try { await api.deleteDomain(id); ElMessage.success('已刪除'); load() } catch (e) { ElMessage.error(e.message) }
}

async function build(id) {
  try { await api.buildSite(id); ElMessage.success('構建完成'); load() } catch (e) { ElMessage.error(e.message) }
=======
  try { await api.deleteDomain(id); ElMessage.success('已删除'); load() } catch (e) { ElMessage.error(e.message) }
}

async function build(id) {
  try { await api.buildSite(id); ElMessage.success('构建完成'); load() } catch (e) { ElMessage.error(e.message) }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
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

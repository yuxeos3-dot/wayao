<template>
  <div>
<<<<<<< HEAD
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
      <h2 style="margin:0">關鍵詞庫</h2>
      <el-upload :http-request="importFile" :show-file-list="false" accept=".csv">
        <el-button type="primary">匯入CSV</el-button>
      </el-upload>
    </div>

    <el-row :gutter="12" style="margin-bottom:12px">
      <el-col :span="4">
        <el-select v-model="filter.category" placeholder="分類" clearable @change="load">
          <el-option v-for="c in categories" :key="c.category" :label="`${c.category} (${c.count})`" :value="c.category" />
        </el-select>
      </el-col>
      <el-col :span="4">
        <el-select v-model="filter.assigned" placeholder="分配狀態" clearable @change="load">
          <el-option label="已分配" value="true" /><el-option label="未分配" value="false" />
        </el-select>
      </el-col>
      <el-col :span="6">
        <el-input v-model="filter.search" placeholder="搜索關鍵詞..." clearable @clear="load" @keyup.enter="load" />
      </el-col>
    </el-row>

    <el-table :data="list" v-loading="loading" size="small">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="keyword" label="關鍵詞" min-width="200" />
      <el-table-column prop="category" label="分類" width="100" />
      <el-table-column prop="monthly_vol" label="搜索量" width="80" sortable />
      <el-table-column prop="difficulty" label="KD" width="60" />
      <el-table-column prop="cpc" label="CPC" width="70" />
      <el-table-column prop="status" label="分配" width="60">
        <template #default="{row}"><el-tag :type="row.status==='assigned'?'success':'info'" size="small">{{ row.status==='assigned'?'是':'否' }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{row}">
          <el-popconfirm title="確定刪除?" @confirm="del(row.id)">
            <template #reference><el-button size="small" type="danger" text>刪除</el-button></template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">关键词库</h1>
        <p class="page-subtitle">管理 SEO 关键词, 支持 CSV 批量导入</p>
      </div>
      <el-upload :http-request="importFile" :show-file-list="false" accept=".csv">
        <el-button type="primary">导入 CSV</el-button>
      </el-upload>
    </div>

    <div class="filter-bar">
      <el-select v-model="filter.category" placeholder="分类" clearable @change="load" style="width:180px">
        <el-option v-for="c in categories" :key="c.category" :label="`${c.category} (${c.count})`" :value="c.category" />
      </el-select>
      <el-select v-model="filter.assigned" placeholder="分配状态" clearable @change="load" style="width:140px">
        <el-option label="已分配" value="true" />
        <el-option label="未分配" value="false" />
      </el-select>
      <el-input v-model="filter.search" placeholder="搜索关键词..." clearable @clear="load" @keyup.enter="load" style="width:220px" />
    </div>

    <div class="section-card">
      <div class="section-body" style="padding:0">
        <el-table :data="list" v-loading="loading" size="small">
          <el-table-column prop="id" label="ID" width="60" />
          <el-table-column prop="keyword" label="关键词" min-width="200">
            <template #default="{row}"><span style="font-weight:500">{{ row.keyword }}</span></template>
          </el-table-column>
          <el-table-column prop="category" label="分类" width="100" />
          <el-table-column prop="monthly_vol" label="搜索量" width="90" sortable />
          <el-table-column prop="difficulty" label="难度" width="70" />
          <el-table-column prop="cpc" label="CPC" width="70" />
          <el-table-column prop="status" label="状态" width="80">
            <template #default="{row}"><el-tag :type="row.status==='assigned'?'success':'info'" size="small">{{ row.status==='assigned'?'已分配':'未分配' }}</el-tag></template>
          </el-table-column>
          <el-table-column label="操作" width="80">
            <template #default="{row}">
              <el-popconfirm title="确定删除?" @confirm="del(row.id)">
                <template #reference><el-button size="small" type="danger" text>删除</el-button></template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31

    <el-pagination style="margin-top:16px;justify-content:flex-end" v-model:current-page="page" :page-size="50" :total="total" @current-change="load" layout="total, prev, pager, next" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const list = ref([])
const loading = ref(false)
const page = ref(1)
const total = ref(0)
const categories = ref([])
const filter = ref({ category: '', assigned: '', search: '' })

async function load() {
  loading.value = true
  try {
    const res = await api.listKeywords({ ...filter.value, page: page.value, size: 50 })
    list.value = res.data || []
    total.value = res.total || 0
  } finally { loading.value = false }
}

async function loadCategories() {
<<<<<<< HEAD
  try {
    const res = await api.keywordCategories()
    categories.value = res.data || []
  } catch {}
=======
  try { const res = await api.keywordCategories(); categories.value = res.data || [] } catch {}
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
}

async function importFile({ file }) {
  const fd = new FormData()
  fd.append('file', file)
  fd.append('market', 'zh-TW')
  try {
    const res = await api.importKeywords(fd)
<<<<<<< HEAD
    ElMessage.success(`匯入 ${res.data?.inserted || 0} 條`)
=======
    ElMessage.success(`成功导入 ${res.data?.inserted || 0} 条`)
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
    load(); loadCategories()
  } catch (e) { ElMessage.error(e.message) }
}

async function del(id) {
  try { await api.deleteKeyword(id); load() } catch (e) { ElMessage.error(e.message) }
}

onMounted(() => { load(); loadCategories() })
</script>

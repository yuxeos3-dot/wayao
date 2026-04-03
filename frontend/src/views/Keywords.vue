<template>
  <div>
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
      <el-table-column prop="volume" label="搜索量" width="80" sortable />
      <el-table-column prop="kd" label="KD" width="60" />
      <el-table-column prop="cpc" label="CPC" width="70" />
      <el-table-column prop="is_assigned" label="分配" width="60">
        <template #default="{row}"><el-tag :type="row.is_assigned?'success':'info'" size="small">{{ row.is_assigned?'是':'否' }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="{row}">
          <el-popconfirm title="確定刪除?" @confirm="del(row.id)">
            <template #reference><el-button size="small" type="danger" text>刪除</el-button></template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

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
  try {
    const res = await api.keywordCategories()
    categories.value = res.data || []
  } catch {}
}

async function importFile({ file }) {
  const fd = new FormData()
  fd.append('file', file)
  fd.append('market', 'zh-TW')
  try {
    const res = await api.importKeywords(fd)
    ElMessage.success(`匯入 ${res.data?.inserted || 0} 條`)
    load(); loadCategories()
  } catch (e) { ElMessage.error(e.message) }
}

async function del(id) {
  try { await api.deleteKeyword(id); load() } catch (e) { ElMessage.error(e.message) }
}

onMounted(() => { load(); loadCategories() })
</script>

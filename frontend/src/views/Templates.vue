<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
      <h2 style="margin:0">模版管理</h2>
      <el-button type="primary" @click="showAdd=true">新增模版</el-button>
    </div>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名稱" />
      <el-table-column prop="slug" label="Slug" />
      <el-table-column prop="css_prefix" label="CSS前綴" width="100" />
      <el-table-column prop="is_active" label="啟用" width="80">
        <template #default="{row}"><el-tag :type="row.is_active?'success':'info'">{{ row.is_active?'是':'否' }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="150">
        <template #default="{row}">
          <el-button size="small" @click="editRow(row)">編輯</el-button>
          <el-popconfirm title="確定刪除?" @confirm="del(row.id)">
            <template #reference><el-button size="small" type="danger">刪除</el-button></template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showAdd" :title="form.id?'編輯模版':'新增模版'" width="500">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名稱"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="Slug"><el-input v-model="form.slug" /></el-form-item>
        <el-form-item label="說明"><el-input v-model="form.description" type="textarea" /></el-form-item>
        <el-form-item label="CSS前綴"><el-input v-model="form.css_prefix" /></el-form-item>
        <el-form-item label="模版路徑"><el-input v-model="form.path" /></el-form-item>
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
const form = ref({ name: '', slug: '', description: '', css_prefix: '', path: '' })

async function load() {
  loading.value = true
  try {
    const res = await api.listTemplates()
    list.value = res.data || []
  } finally { loading.value = false }
}

function editRow(row) {
  form.value = { ...row }
  showAdd.value = true
}

async function save() {
  try {
    if (form.value.id) {
      await api.updateTemplate(form.value.id, form.value)
    } else {
      await api.createTemplate(form.value)
    }
    ElMessage.success('已保存')
    showAdd.value = false
    form.value = { name: '', slug: '', description: '', css_prefix: '', path: '' }
    load()
  } catch (e) { ElMessage.error(e.message) }
}

async function del(id) {
  try {
    await api.deleteTemplate(id)
    ElMessage.success('已刪除')
    load()
  } catch (e) { ElMessage.error(e.message) }
}

onMounted(load)
</script>

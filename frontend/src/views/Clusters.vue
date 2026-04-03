<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
      <h2 style="margin:0">站群管理</h2>
      <el-button type="primary" @click="showAdd=true">新增站群</el-button>
    </div>

    <el-table :data="list" v-loading="loading">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" label="名稱" />
      <el-table-column prop="strategy" label="策略" width="100" />
      <el-table-column prop="member_count" label="成員數" width="80" />
      <el-table-column prop="created_at" label="建立時間" width="170" />
      <el-table-column label="操作" width="120">
        <template #default="{row}">
          <el-popconfirm title="確定刪除?" @confirm="del(row.id)">
            <template #reference><el-button size="small" type="danger">刪除</el-button></template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showAdd" title="新增站群" width="400">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名稱"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="策略">
          <el-select v-model="form.strategy"><el-option label="輪鏈 (wheel)" value="wheel" /><el-option label="金字塔 (pyramid)" value="pyramid" /><el-option label="隨機 (random)" value="random" /></el-select>
        </el-form-item>
      </el-form>
      <template #footer><el-button @click="showAdd=false">取消</el-button><el-button type="primary" @click="save">保存</el-button></template>
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
const form = ref({ name: '', strategy: 'wheel' })

async function load() {
  loading.value = true
  try { const res = await api.listClusters(); list.value = res.data || [] } finally { loading.value = false }
}
async function save() {
  try { await api.createCluster(form.value); ElMessage.success('已建立'); showAdd.value = false; form.value = { name: '', strategy: 'wheel' }; load() } catch (e) { ElMessage.error(e.message) }
}
async function del(id) {
  try { await api.deleteCluster(id); ElMessage.success('已刪除'); load() } catch (e) { ElMessage.error(e.message) }
}
onMounted(load)
</script>

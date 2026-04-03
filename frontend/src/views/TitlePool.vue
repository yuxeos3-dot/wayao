<template>
  <div>
    <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:16px">
      <h2 style="margin:0">標題變體池</h2>
      <el-button type="primary" @click="showAdd=true">新增變體</el-button>
    </div>

    <el-select v-model="filter" placeholder="篩選詞類" clearable @change="load" style="margin-bottom:12px">
      <el-option v-for="t in kwTypes" :key="t" :label="t" :value="t" />
    </el-select>

    <el-table :data="list" v-loading="loading" size="small">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="keyword_type" label="詞類" width="100" />
      <el-table-column prop="pattern" label="標題模式" min-width="300">
        <template #default="{row}"><code>{{ row.pattern }}</code></template>
      </el-table-column>
      <el-table-column prop="is_active" label="啟用" width="60">
        <template #default="{row}"><el-tag :type="row.is_active?'success':'info'" size="small">{{ row.is_active?'是':'否' }}</el-tag></template>
      </el-table-column>
      <el-table-column label="操作" width="80">
        <template #default="{row}">
          <el-popconfirm title="確定刪除?" @confirm="del(row.id)">
            <template #reference><el-button size="small" type="danger" text>刪除</el-button></template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showAdd" title="新增標題變體" width="500">
      <el-form :model="form" label-width="80px">
        <el-form-item label="詞類"><el-select v-model="form.keyword_type"><el-option v-for="t in kwTypes" :key="t" :label="t" :value="t" /></el-select></el-form-item>
        <el-form-item label="標題模式"><el-input v-model="form.pattern" placeholder="{keyword} - {year}最新評測｜專業推薦" /><div style="color:#909399;font-size:12px;margin-top:4px">可用變數: {keyword}, {year}</div></el-form-item>
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
const filter = ref('')
const form = ref({ keyword_type: 'brand', pattern: '' })
const kwTypes = ['brand','game','sports','generic','promo','payment','affiliate','strategy','app','register','region','credit','live','community','terms']

async function load() {
  loading.value = true
  try { const res = await api.listTitlePool({ keyword_type: filter.value }); list.value = res.data || [] } finally { loading.value = false }
}
async function save() {
  try { await api.createTitleVariant(form.value); ElMessage.success('已新增'); showAdd.value = false; form.value = { keyword_type: 'brand', pattern: '' }; load() } catch (e) { ElMessage.error(e.message) }
}
async function del(id) {
  try { await api.deleteTitleVariant(id); load() } catch (e) { ElMessage.error(e.message) }
}
onMounted(load)
</script>

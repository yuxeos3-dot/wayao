<template>
  <div>
    <h2>城市矩陣</h2>
    <el-row :gutter="12" style="margin-bottom:12px">
      <el-col :span="6">
        <el-select v-model="domainId" placeholder="選擇域名" filterable @change="loadMatrix">
          <el-option v-for="d in domains" :key="d.id" :label="d.domain" :value="d.id" />
        </el-select>
      </el-col>
      <el-col :span="4"><el-button type="primary" @click="addCity" :disabled="!domainId">新增城市</el-button></el-col>
    </el-row>

    <el-table :data="cities" v-loading="loading" size="small" v-if="domainId">
      <el-table-column label="城市名" width="120">
        <template #default="{row,$index}"><el-input v-model="cities[$index].city_name" size="small" /></template>
      </el-table-column>
      <el-table-column label="Slug" width="120">
        <template #default="{row,$index}"><el-input v-model="cities[$index].city_slug" size="small" /></template>
      </el-table-column>
      <el-table-column label="額外標題" min-width="200">
        <template #default="{row,$index}"><el-input v-model="cities[$index].extra_title" size="small" /></template>
      </el-table-column>
      <el-table-column label="額外描述" min-width="200">
        <template #default="{row,$index}"><el-input v-model="cities[$index].extra_desc" size="small" /></template>
      </el-table-column>
      <el-table-column label="操作" width="60">
        <template #default="{$index}"><el-button size="small" type="danger" text @click="cities.splice($index,1)">刪</el-button></template>
      </el-table-column>
    </el-table>

    <el-button type="primary" style="margin-top:16px" @click="saveMatrix" :disabled="!domainId" :loading="saving">保存</el-button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const domains = ref([])
const domainId = ref('')
const cities = ref([])
const loading = ref(false)
const saving = ref(false)

function addCity() { cities.value.push({ city_name: '', city_slug: '', extra_title: '', extra_desc: '' }) }

async function loadMatrix() {
  if (!domainId.value) return
  loading.value = true
  try { const res = await api.getCityMatrix(domainId.value); cities.value = res.data || [] } finally { loading.value = false }
}

async function saveMatrix() {
  saving.value = true
  try { await api.saveCityMatrix(domainId.value, { cities: cities.value }); ElMessage.success('已保存') } catch (e) { ElMessage.error(e.message) } finally { saving.value = false }
}

onMounted(async () => {
  try { const res = await api.listDomains({ size: 500 }); domains.value = res.data || [] } catch {}
})
</script>

<template>
  <div>
<<<<<<< HEAD
    <h2>收錄監控</h2>
    <el-row :gutter="12" style="margin-bottom:12px">
      <el-col :span="6">
        <el-select v-model="domainId" placeholder="選擇域名" filterable @change="checkIndex">
          <el-option v-for="d in domains" :key="d.id" :label="d.domain" :value="d.id" />
        </el-select>
      </el-col>
      <el-col :span="4"><el-button type="primary" @click="batchCheck" :loading="loading">批量檢查</el-button></el-col>
    </el-row>

    <el-card v-if="indexResult">
      <template #header>收錄狀態</template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="域名">{{ indexResult.domain }}</el-descriptions-item>
        <el-descriptions-item label="狀態">{{ indexResult.status }}</el-descriptions-item>
      </el-descriptions>
      <el-alert type="info" style="margin-top:12px" :closable="false">
        {{ indexResult.tip }}
      </el-alert>
    </el-card>

    <el-card style="margin-top:16px">
      <template #header>IndexNow 推送記錄</template>
      <el-table :data="submissions" size="small" v-if="domainId">
        <el-table-column prop="url" label="URL" min-width="250" />
        <el-table-column prop="engine" label="引擎" width="100" />
        <el-table-column prop="status" label="狀態" width="80">
          <template #default="{row}"><el-tag :type="row.status==='submitted'?'success':'info'" size="small">{{ row.status }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" label="時間" width="170" />
      </el-table>
      <el-empty v-else description="選擇域名查看推送記錄" />
    </el-card>
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">收录监控</h1>
        <p class="page-subtitle">搜索引擎收录状态与 IndexNow 推送记录</p>
      </div>
      <el-button type="primary" @click="batchCheck" :loading="loading">批量检查</el-button>
    </div>

    <div class="filter-bar">
      <el-select v-model="domainId" placeholder="选择域名" filterable clearable @change="checkIndex" style="width:280px">
        <el-option v-for="d in domains" :key="d.id" :label="d.domain" :value="d.id" />
      </el-select>
    </div>

    <el-row :gutter="20">
      <el-col :span="indexResult ? 12 : 24">
        <div class="section-card">
          <div class="section-header">IndexNow 推送记录</div>
          <div class="section-body" :style="domainId && submissions.length ? {padding:0} : {}">
            <el-table :data="submissions" size="small" v-if="domainId && submissions.length">
              <el-table-column prop="url" label="URL" min-width="250" />
              <el-table-column prop="engine" label="引擎" width="100" />
              <el-table-column prop="status" label="状态" width="80">
                <template #default="{row}"><el-tag :type="row.status==='submitted'?'success':'info'" size="small">{{ row.status==='submitted'?'已推送':row.status }}</el-tag></template>
              </el-table-column>
              <el-table-column prop="created_at" label="时间" width="170" />
            </el-table>
            <el-empty v-else :description="domainId ? '暂无推送记录' : '请选择域名查看'" :image-size="60" />
          </div>
        </div>
      </el-col>
      <el-col :span="12" v-if="indexResult">
        <div class="section-card">
          <div class="section-header">收录状态</div>
          <div class="section-body">
            <el-descriptions :column="1" border>
              <el-descriptions-item label="域名">{{ indexResult.domain }}</el-descriptions-item>
              <el-descriptions-item label="状态">{{ indexResult.status }}</el-descriptions-item>
            </el-descriptions>
            <el-alert type="info" style="margin-top:16px" :closable="false">{{ indexResult.tip }}</el-alert>
          </div>
        </div>
      </el-col>
    </el-row>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const domains = ref([])
const domainId = ref('')
const indexResult = ref(null)
const submissions = ref([])
const loading = ref(false)

async function checkIndex() {
<<<<<<< HEAD
  if (!domainId.value) return
=======
  if (!domainId.value) { indexResult.value = null; submissions.value = []; return }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
  try {
    const res = await api.get(`/index-status/${domainId.value}`)
    indexResult.value = res.data
    const sub = await api.getIndexNowRecords(domainId.value)
    submissions.value = sub.data || []
  } catch (e) { ElMessage.error(e.message) }
}

async function batchCheck() {
  loading.value = true
<<<<<<< HEAD
  try {
    await api.post('/index-status/batch')
    ElMessage.success('批量檢查已啟動')
  } catch (e) { ElMessage.error(e.message) }
=======
  try { await api.post('/index-status/batch'); ElMessage.success('批量检查已启动') } catch (e) { ElMessage.error(e.message) }
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
  finally { loading.value = false }
}

onMounted(async () => {
  try { const res = await api.listDomains({ size: 500 }); domains.value = res.data || [] } catch {}
})
</script>

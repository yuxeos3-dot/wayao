<template>
  <div>
    <h2>安全防護</h2>
    <el-tabs>
      <el-tab-pane label="IP規則">
        <div style="margin-bottom:12px">
          <el-input v-model="ipForm.ip" placeholder="IP地址" style="width:200px;margin-right:8px" />
          <el-select v-model="ipForm.action" style="width:100px;margin-right:8px"><el-option label="封鎖" value="block" /><el-option label="允許" value="allow" /></el-select>
          <el-input v-model="ipForm.reason" placeholder="原因" style="width:200px;margin-right:8px" />
          <el-button type="primary" @click="addIP">新增</el-button>
        </div>
        <el-table :data="ipRules" size="small">
          <el-table-column prop="ip" label="IP" />
          <el-table-column prop="action" label="動作" width="80">
            <template #default="{row}"><el-tag :type="row.action==='block'?'danger':'success'" size="small">{{ row.action }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" />
          <el-table-column prop="created_at" label="時間" width="170" />
          <el-table-column label="操作" width="80">
            <template #default="{row}">
              <el-popconfirm title="確定刪除?" @confirm="delIP(row.id)"><template #reference><el-button size="small" type="danger" text>刪除</el-button></template></el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="UA規則">
        <div style="margin-bottom:12px">
          <el-input v-model="uaForm.pattern" placeholder="UA關鍵詞" style="width:200px;margin-right:8px" />
          <el-select v-model="uaForm.action" style="width:100px;margin-right:8px"><el-option label="封鎖" value="block" /><el-option label="允許" value="allow" /></el-select>
          <el-button type="primary" @click="addUA">新增</el-button>
        </div>
        <el-table :data="uaRules" size="small">
          <el-table-column prop="pattern" label="Pattern" />
          <el-table-column prop="rule_type" label="動作" width="80">
            <template #default="{row}"><el-tag :type="row.rule_type==='block'?'danger':'success'" size="small">{{ row.rule_type }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="is_active" label="啟用" width="60">
            <template #default="{row}"><el-tag :type="row.is_active?'success':'info'" size="small">{{ row.is_active?'是':'否' }}</el-tag></template>
          </el-table-column>
          <el-table-column prop="created_at" label="時間" width="170" />
          <el-table-column label="操作" width="80">
            <template #default="{row}">
              <el-popconfirm title="確定刪除?" @confirm="delUA(row.id)"><template #reference><el-button size="small" type="danger" text>刪除</el-button></template></el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

const ipRules = ref([])
const uaRules = ref([])
const ipForm = ref({ ip: '', action: 'block', reason: '' })
const uaForm = ref({ pattern: '', action: 'block' })

async function loadIP() { try { const res = await api.listIPRules(); ipRules.value = res.data || [] } catch {} }
async function loadUA() { try { const res = await api.listUARules(); uaRules.value = res.data || [] } catch {} }

async function addIP() {
  if (!ipForm.value.ip) return
  try { await api.addIPRule(ipForm.value); ElMessage.success('已新增'); ipForm.value = { ip: '', action: 'block', reason: '' }; loadIP() } catch (e) { ElMessage.error(e.message) }
}
async function addUA() {
  if (!uaForm.value.pattern) return
  try { await api.addUARule(uaForm.value); ElMessage.success('已新增'); uaForm.value = { pattern: '', action: 'block' }; loadUA() } catch (e) { ElMessage.error(e.message) }
}
async function delIP(id) { try { await api.deleteIPRule(id); loadIP() } catch {} }
async function delUA(id) { try { await api.deleteUARule(id); loadUA() } catch {} }

onMounted(() => { loadIP(); loadUA() })
</script>

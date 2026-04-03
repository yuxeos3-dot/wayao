<template>
  <div>
<<<<<<< HEAD
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
=======
    <div class="page-header">
      <div>
        <h1 class="page-title">安全防护</h1>
        <p class="page-subtitle">IP 封禁规则与 User-Agent 过滤规则管理</p>
      </div>
    </div>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="IP 规则" name="ip">
        <div class="section-card">
          <div class="section-header">
            <span>IP 规则列表</span>
            <div style="display:flex;gap:8px;align-items:center">
              <el-input v-model="ipForm.ip" placeholder="IP 地址" style="width:180px" size="small" />
              <el-select v-model="ipForm.action" style="width:100px" size="small">
                <el-option label="封禁" value="block" />
                <el-option label="放行" value="allow" />
              </el-select>
              <el-input v-model="ipForm.reason" placeholder="原因" style="width:160px" size="small" />
              <el-button type="primary" size="small" @click="addIP">新增</el-button>
            </div>
          </div>
          <div class="section-body" style="padding:0">
            <el-table :data="ipRules" size="small">
              <el-table-column prop="ip" label="IP 地址" />
              <el-table-column prop="action" label="动作" width="80">
                <template #default="{row}"><el-tag :type="row.action==='block'?'danger':'success'" size="small">{{ row.action==='block'?'封禁':'放行' }}</el-tag></template>
              </el-table-column>
              <el-table-column prop="reason" label="原因" />
              <el-table-column prop="created_at" label="创建时间" width="170" />
              <el-table-column label="" width="80">
                <template #default="{row}">
                  <el-popconfirm title="确定删除?" @confirm="delIP(row.id)">
                    <template #reference><el-button size="small" type="danger" text>删除</el-button></template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="UA 规则" name="ua">
        <div class="section-card">
          <div class="section-header">
            <span>UA 规则列表</span>
            <div style="display:flex;gap:8px;align-items:center">
              <el-input v-model="uaForm.pattern" placeholder="UA 关键词" style="width:180px" size="small" />
              <el-select v-model="uaForm.action" style="width:100px" size="small">
                <el-option label="封禁" value="block" />
                <el-option label="放行" value="allow" />
              </el-select>
              <el-button type="primary" size="small" @click="addUA">新增</el-button>
            </div>
          </div>
          <div class="section-body" style="padding:0">
            <el-table :data="uaRules" size="small">
              <el-table-column prop="pattern" label="匹配模式" />
              <el-table-column prop="rule_type" label="动作" width="80">
                <template #default="{row}"><el-tag :type="row.rule_type==='block'?'danger':'success'" size="small">{{ row.rule_type==='block'?'封禁':'放行' }}</el-tag></template>
              </el-table-column>
              <el-table-column prop="is_active" label="状态" width="80">
                <template #default="{row}"><el-tag :type="row.is_active?'success':'info'" size="small">{{ row.is_active?'启用':'停用' }}</el-tag></template>
              </el-table-column>
              <el-table-column prop="created_at" label="创建时间" width="170" />
              <el-table-column label="" width="80">
                <template #default="{row}">
                  <el-popconfirm title="确定删除?" @confirm="delUA(row.id)">
                    <template #reference><el-button size="small" type="danger" text>删除</el-button></template>
                  </el-popconfirm>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'
import { ElMessage } from 'element-plus'

<<<<<<< HEAD
=======
const activeTab = ref('ip')
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
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

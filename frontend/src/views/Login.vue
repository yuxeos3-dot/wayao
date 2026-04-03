<template>
  <div style="display:flex;justify-content:center;align-items:center;height:100vh;background:#f0f2f5">
    <el-card style="width:400px">
      <template #header><h2 style="text-align:center;margin:0">BrandSite Pro</h2></template>
      <el-form @submit.prevent="handleLogin">
        <el-form-item label="密碼">
          <el-input v-model="password" type="password" placeholder="首次登入密碼為 admin" show-password @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" style="width:100%" @click="handleLogin">登入</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const password = ref('')
const loading = ref(false)
const auth = useAuthStore()

async function handleLogin() {
  if (!password.value) return
  loading.value = true
  try {
    await auth.login(password.value)
    ElMessage.success('登入成功')
  } catch (e) {
    ElMessage.error(e.message || '登入失敗')
  } finally {
    loading.value = false
  }
}
</script>

<template>
<<<<<<< HEAD
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
=======
  <div class="login-page">
    <div class="login-container">
      <div class="login-brand">
        <div class="login-icon">W</div>
        <h1 class="login-title">Wayao CMS</h1>
        <p class="login-subtitle">SEO 静态站群管理系统</p>
      </div>

      <form @submit.prevent="handleLogin" class="login-form">
        <div class="input-group">
          <label class="input-label">登录密码</label>
          <el-input
            v-model="password"
            type="password"
            placeholder="首次登录默认密码 admin"
            show-password
            size="large"
            @keyup.enter="handleLogin"
          />
        </div>
        <el-button
          type="primary"
          size="large"
          :loading="loading"
          style="width:100%;height:48px;font-size:16px;font-weight:600"
          @click="handleLogin"
        >
          登录
        </el-button>
      </form>

      <p class="login-footer">Wayao CMS v2.0</p>
    </div>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
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
<<<<<<< HEAD
    ElMessage.success('登入成功')
  } catch (e) {
    ElMessage.error(e.message || '登入失敗')
=======
    ElMessage.success('登录成功')
  } catch (e) {
    ElMessage.error(e.message || '登录失败')
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
  } finally {
    loading.value = false
  }
}
</script>
<<<<<<< HEAD
=======

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f5f7 0%, #e8e8ed 100%);
}

.login-container {
  width: 400px;
  background: var(--wayao-card);
  border-radius: 20px;
  padding: 48px 40px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.08), 0 1px 3px rgba(0,0,0,0.04);
}

.login-brand {
  text-align: center;
  margin-bottom: 40px;
}

.login-icon {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  background: linear-gradient(135deg, var(--wayao-accent), #5ac8fa);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 700;
  font-size: 32px;
  letter-spacing: -0.02em;
  margin-bottom: 16px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: var(--wayao-text);
  margin: 0;
}

.login-subtitle {
  font-size: 14px;
  color: var(--wayao-text-secondary);
  margin-top: 8px;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-group { display: flex; flex-direction: column; gap: 8px; }
.input-label { font-size: 14px; font-weight: 600; color: var(--wayao-text); }

.login-footer {
  text-align: center;
  margin-top: 32px;
  font-size: 12px;
  color: var(--wayao-text-tertiary);
}
</style>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31

<template>
  <el-container style="height:100vh">
    <el-aside width="220px" style="background:#304156">
      <div style="padding:20px;text-align:center;color:#fff;font-size:18px;font-weight:bold">BrandSite Pro</div>
      <el-menu :default-active="$route.path" background-color="#304156" text-color="#bfcbd9" active-text-color="#409EFF" router>
        <el-menu-item index="/dashboard"><el-icon><DataBoard /></el-icon><span>數據看板</span></el-menu-item>
        <el-menu-item index="/templates"><el-icon><Document /></el-icon><span>模版管理</span></el-menu-item>
        <el-menu-item index="/domains"><el-icon><Link /></el-icon><span>域名管理</span></el-menu-item>
        <el-menu-item index="/keywords"><el-icon><Key /></el-icon><span>關鍵詞庫</span></el-menu-item>
        <el-menu-item index="/builder"><el-icon><Cpu /></el-icon><span>構建部署</span></el-menu-item>
        <el-menu-item index="/rankings"><el-icon><TrendCharts /></el-icon><span>排名追蹤</span></el-menu-item>
        <el-sub-menu index="v4">
          <template #title><el-icon><Tools /></el-icon><span>進階功能</span></template>
          <el-menu-item index="/title-pool">標題變體池</el-menu-item>
          <el-menu-item index="/city-matrix">城市矩陣</el-menu-item>
          <el-menu-item index="/clusters">站群管理</el-menu-item>
          <el-menu-item index="/health-check">健康檢查</el-menu-item>
          <el-menu-item index="/index-status">收錄監控</el-menu-item>
          <el-menu-item index="/security">安全防護</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/settings"><el-icon><Setting /></el-icon><span>系統設定</span></el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="display:flex;justify-content:flex-end;align-items:center;border-bottom:1px solid #e6e6e6">
        <el-button type="danger" text @click="auth.logout()">登出</el-button>
      </el-header>
      <el-main style="background:#f0f2f5;overflow-y:auto">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
const auth = useAuthStore()
</script>

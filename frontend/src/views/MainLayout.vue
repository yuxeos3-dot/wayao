<template>
<<<<<<< HEAD
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
=======
  <div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-brand">
        <div class="brand-icon">W</div>
        <span class="brand-text">Wayao CMS</span>
      </div>

      <nav class="sidebar-nav">
        <div class="nav-section">
          <router-link to="/dashboard" class="nav-item" :class="{active: isActive('/dashboard')}">
            <el-icon><DataBoard /></el-icon><span>数据看板</span>
          </router-link>
        </div>

        <div class="nav-section">
          <div class="nav-label">站点管理</div>
          <router-link to="/sites" class="nav-item" :class="{active: isActive('/sites')}">
            <el-icon><Link /></el-icon><span>域名管理</span>
          </router-link>
          <router-link to="/templates" class="nav-item" :class="{active: isActive('/templates')}">
            <el-icon><Document /></el-icon><span>模版管理</span>
          </router-link>
          <router-link to="/clusters" class="nav-item" :class="{active: isActive('/clusters')}">
            <el-icon><Grid /></el-icon><span>站群管理</span>
          </router-link>
        </div>

        <div class="nav-section">
          <div class="nav-label">内容中心</div>
          <router-link to="/keywords" class="nav-item" :class="{active: isActive('/keywords')}">
            <el-icon><Key /></el-icon><span>关键词库</span>
          </router-link>
          <router-link to="/title-pool" class="nav-item" :class="{active: isActive('/title-pool')}">
            <el-icon><Collection /></el-icon><span>标题池</span>
          </router-link>
          <router-link to="/city-matrix" class="nav-item" :class="{active: isActive('/city-matrix')}">
            <el-icon><MapLocation /></el-icon><span>城市矩阵</span>
          </router-link>
        </div>

        <div class="nav-section">
          <div class="nav-label">构建发布</div>
          <router-link to="/builder" class="nav-item" :class="{active: isActive('/builder')}">
            <el-icon><Cpu /></el-icon><span>构建部署</span>
          </router-link>
          <router-link to="/index-status" class="nav-item" :class="{active: isActive('/index-status')}">
            <el-icon><Search /></el-icon><span>收录监控</span>
          </router-link>
        </div>

        <div class="nav-section">
          <div class="nav-label">数据分析</div>
          <router-link to="/rankings" class="nav-item" :class="{active: isActive('/rankings')}">
            <el-icon><TrendCharts /></el-icon><span>排名追踪</span>
          </router-link>
          <router-link to="/health-check" class="nav-item" :class="{active: isActive('/health-check')}">
            <el-icon><CircleCheck /></el-icon><span>健康检查</span>
          </router-link>
        </div>
      </nav>

      <div class="sidebar-footer">
        <router-link to="/settings" class="nav-item" :class="{active: isActive('/settings')}">
          <el-icon><Setting /></el-icon><span>系统设置</span>
        </router-link>
        <router-link to="/security" class="nav-item" :class="{active: isActive('/security')}">
          <el-icon><Lock /></el-icon><span>安全防护</span>
        </router-link>
        <div class="nav-item logout" @click="auth.logout()">
          <el-icon><SwitchButton /></el-icon><span>退出登录</span>
        </div>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const auth = useAuthStore()

function isActive(path) {
  return route.path === path || route.path.startsWith(path + '/')
}
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  width: 240px;
  min-width: 240px;
  background: var(--wayao-sidebar);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 24px 20px;
  border-bottom: 1px solid rgba(255,255,255,0.06);
}

.brand-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--wayao-accent), #5ac8fa);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-weight: 700;
  font-size: 18px;
  letter-spacing: -0.02em;
  flex-shrink: 0;
}

.brand-text {
  color: #fff;
  font-size: 17px;
  font-weight: 600;
  letter-spacing: -0.02em;
}

.sidebar-nav {
  flex: 1;
  padding: 8px 0;
  overflow-y: auto;
}

.nav-section {
  padding: 4px 0;
}

.nav-label {
  padding: 16px 20px 6px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: rgba(255,255,255,0.35);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 16px;
  margin: 1px 8px;
  border-radius: 8px;
  color: rgba(255,255,255,0.65);
  font-size: 14px;
  font-weight: 450;
  text-decoration: none;
  cursor: pointer;
  transition: var(--wayao-transition);
}
.nav-item:hover {
  background: var(--wayao-sidebar-hover);
  color: rgba(255,255,255,0.9);
}
.nav-item.active {
  background: var(--wayao-sidebar-active);
  color: #fff;
  font-weight: 500;
}
.nav-item .el-icon { font-size: 18px; opacity: 0.8; }
.nav-item.active .el-icon { opacity: 1; }

.sidebar-footer {
  padding: 8px 0 16px;
  border-top: 1px solid rgba(255,255,255,0.06);
}

.nav-item.logout { color: rgba(255,59,48,0.8); }
.nav-item.logout:hover { background: rgba(255,59,48,0.1); color: #ff3b30; }

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 32px 40px;
  background: var(--wayao-bg);
}
</style>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31

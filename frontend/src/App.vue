<template>
  <router-view />
</template>
<<<<<<< HEAD
=======

<style>
:root {
  --wayao-bg: #f5f5f7;
  --wayao-card: #ffffff;
  --wayao-sidebar: #1d1d1f;
  --wayao-sidebar-hover: rgba(255,255,255,0.08);
  --wayao-sidebar-active: rgba(255,255,255,0.12);
  --wayao-text: #1d1d1f;
  --wayao-text-secondary: #86868b;
  --wayao-text-tertiary: #aeaeb2;
  --wayao-border: rgba(0,0,0,0.06);
  --wayao-accent: #0071e3;
  --wayao-accent-hover: #0077ed;
  --wayao-success: #34c759;
  --wayao-warning: #ff9f0a;
  --wayao-danger: #ff3b30;
  --wayao-radius: 12px;
  --wayao-radius-sm: 8px;
  --wayao-shadow: 0 1px 3px rgba(0,0,0,0.04), 0 1px 2px rgba(0,0,0,0.06);
  --wayao-shadow-md: 0 4px 12px rgba(0,0,0,0.08);
  --wayao-shadow-lg: 0 8px 30px rgba(0,0,0,0.12);
  --wayao-transition: all 0.2s cubic-bezier(0.25, 0.1, 0.25, 1);
}

* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'SF Pro Display', 'SF Pro Text', 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: var(--wayao-text);
  background: var(--wayao-bg);
}

/* === Element Plus Override: Apple-like Design === */

/* Cards */
.el-card {
  border: none !important;
  border-radius: var(--wayao-radius) !important;
  box-shadow: var(--wayao-shadow) !important;
  transition: var(--wayao-transition);
}
.el-card:hover {
  box-shadow: var(--wayao-shadow-md) !important;
}
.el-card__header {
  border-bottom: 1px solid var(--wayao-border) !important;
  padding: 16px 20px !important;
  font-weight: 600 !important;
  font-size: 15px !important;
  color: var(--wayao-text) !important;
}
.el-card__body { padding: 20px !important; }

/* Buttons */
.el-button {
  border-radius: var(--wayao-radius-sm) !important;
  font-weight: 500 !important;
  transition: var(--wayao-transition) !important;
  border: none !important;
  letter-spacing: -0.01em;
}
.el-button--primary {
  background: var(--wayao-accent) !important;
  box-shadow: 0 1px 2px rgba(0,113,227,0.3) !important;
}
.el-button--primary:hover {
  background: var(--wayao-accent-hover) !important;
  box-shadow: 0 2px 8px rgba(0,113,227,0.4) !important;
  transform: translateY(-0.5px);
}
.el-button--success {
  background: var(--wayao-success) !important;
  box-shadow: 0 1px 2px rgba(52,199,89,0.3) !important;
}
.el-button--danger {
  background: var(--wayao-danger) !important;
  box-shadow: 0 1px 2px rgba(255,59,48,0.3) !important;
}
.el-button--default {
  background: rgba(0,0,0,0.04) !important;
  color: var(--wayao-text) !important;
}
.el-button--default:hover {
  background: rgba(0,0,0,0.08) !important;
}

/* Tags */
.el-tag {
  border: none !important;
  border-radius: 6px !important;
  font-weight: 500 !important;
  font-size: 12px !important;
}
.el-tag--success { background: rgba(52,199,89,0.12) !important; color: #248a3d !important; }
.el-tag--warning { background: rgba(255,159,10,0.12) !important; color: #c93400 !important; }
.el-tag--danger { background: rgba(255,59,48,0.12) !important; color: #d70015 !important; }
.el-tag--info { background: rgba(142,142,147,0.12) !important; color: #636366 !important; }
.el-tag--primary { background: rgba(0,113,227,0.1) !important; color: var(--wayao-accent) !important; }

/* Table */
.el-table {
  border-radius: var(--wayao-radius) !important;
  overflow: hidden;
  --el-table-border-color: var(--wayao-border);
  --el-table-header-bg-color: rgba(0,0,0,0.02);
}
.el-table th.el-table__cell {
  font-weight: 600 !important;
  font-size: 12px !important;
  text-transform: uppercase !important;
  letter-spacing: 0.02em !important;
  color: var(--wayao-text-secondary) !important;
}
.el-table td.el-table__cell { font-size: 14px !important; }
.el-table--enable-row-hover .el-table__body tr:hover > td.el-table__cell {
  background: rgba(0,113,227,0.03) !important;
}
.el-table::before { display: none !important; }

/* Input */
.el-input__wrapper {
  border-radius: var(--wayao-radius-sm) !important;
  box-shadow: 0 0 0 1px rgba(0,0,0,0.1) inset !important;
  transition: var(--wayao-transition) !important;
}
.el-input__wrapper:hover { box-shadow: 0 0 0 1px rgba(0,0,0,0.2) inset !important; }
.el-input__wrapper.is-focus { box-shadow: 0 0 0 2px var(--wayao-accent) inset !important; }
.el-textarea__inner {
  border-radius: var(--wayao-radius-sm) !important;
  border-color: rgba(0,0,0,0.1) !important;
}
.el-textarea__inner:focus { border-color: var(--wayao-accent) !important; box-shadow: 0 0 0 1px var(--wayao-accent) !important; }

/* Select */
.el-select .el-input__wrapper { border-radius: var(--wayao-radius-sm) !important; }

/* Dialog */
.el-dialog {
  border-radius: var(--wayao-radius) !important;
  box-shadow: var(--wayao-shadow-lg) !important;
  overflow: hidden;
}
.el-dialog__header {
  border-bottom: 1px solid var(--wayao-border);
  padding: 16px 24px !important;
  margin-right: 0 !important;
}
.el-dialog__title { font-weight: 600 !important; font-size: 16px !important; }
.el-dialog__body { padding: 24px !important; }
.el-dialog__footer { padding: 16px 24px !important; border-top: 1px solid var(--wayao-border); }

/* Tabs */
.el-tabs__item {
  font-weight: 500 !important;
  color: var(--wayao-text-secondary) !important;
  transition: var(--wayao-transition) !important;
}
.el-tabs__item.is-active { color: var(--wayao-accent) !important; font-weight: 600 !important; }
.el-tabs__active-bar { background-color: var(--wayao-accent) !important; border-radius: 2px; }

/* Pagination */
.el-pagination { --el-pagination-hover-color: var(--wayao-accent); }
.el-pager li {
  border-radius: var(--wayao-radius-sm) !important;
  font-weight: 500 !important;
}

/* Switch */
.el-switch.is-checked .el-switch__core { background-color: var(--wayao-accent) !important; border-color: var(--wayao-accent) !important; }

/* Empty */
.el-empty__description p { color: var(--wayao-text-secondary) !important; }

/* Alert */
.el-alert { border-radius: var(--wayao-radius-sm) !important; border: none !important; }

/* Form */
.el-form-item__label { font-weight: 500 !important; color: var(--wayao-text-secondary) !important; font-size: 13px !important; }

/* Descriptions */
.el-descriptions { border-radius: var(--wayao-radius-sm) !important; }

/* Scrollbar */
::-webkit-scrollbar { width: 6px; height: 6px; }
::-webkit-scrollbar-track { background: transparent; }
::-webkit-scrollbar-thumb { background: rgba(0,0,0,0.15); border-radius: 3px; }
::-webkit-scrollbar-thumb:hover { background: rgba(0,0,0,0.25); }

/* Page headers */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.page-title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.02em;
  color: var(--wayao-text);
}
.page-subtitle {
  font-size: 14px;
  color: var(--wayao-text-secondary);
  margin-top: 4px;
}

/* Stat cards */
.stat-grid { display: grid; gap: 16px; margin-bottom: 24px; }
.stat-grid-4 { grid-template-columns: repeat(4, 1fr); }
.stat-card {
  background: var(--wayao-card);
  border-radius: var(--wayao-radius);
  padding: 20px 24px;
  box-shadow: var(--wayao-shadow);
  transition: var(--wayao-transition);
}
.stat-card:hover { box-shadow: var(--wayao-shadow-md); transform: translateY(-1px); }
.stat-label { font-size: 13px; color: var(--wayao-text-secondary); font-weight: 500; margin-bottom: 8px; }
.stat-value { font-size: 32px; font-weight: 700; letter-spacing: -0.03em; color: var(--wayao-text); }
.stat-change { font-size: 12px; font-weight: 600; margin-top: 4px; }
.stat-change.up { color: var(--wayao-success); }
.stat-change.down { color: var(--wayao-danger); }

/* Section card */
.section-card {
  background: var(--wayao-card);
  border-radius: var(--wayao-radius);
  box-shadow: var(--wayao-shadow);
  overflow: hidden;
  margin-bottom: 24px;
}
.section-header {
  padding: 16px 24px;
  border-bottom: 1px solid var(--wayao-border);
  font-size: 15px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.section-body { padding: 24px; }

/* Filter bar */
.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  align-items: center;
}
</style>
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31

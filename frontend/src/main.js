import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
<<<<<<< HEAD
import zhTw from 'element-plus/es/locale/lang/zh-tw'
=======
import zhCn from 'element-plus/es/locale/lang/zh-cn'
>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
<<<<<<< HEAD
app.use(ElementPlus, { locale: zhTw })

// register all icons globally
=======
app.use(ElementPlus, { locale: zhCn })

>>>>>>> 90adefdc839ffaeedc55c4dded5e12b4fcc7ec31
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.mount('#app')

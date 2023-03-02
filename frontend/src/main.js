import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import Axios from 'axios'
export const bus = new Vue()

Axios.defaults.baseURL = 'http://192.168.20.149:8080/'
Vue.config.productionTip = false

new Vue({
  bus,
  vuetify,
  render: h => h(App)
}).$mount('#app')

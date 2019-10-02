import Vue from 'vue'
import routes from './routes'
import VueRouter from 'vue-router'
import App from './App.vue'

import 'normalize.css'

Vue.config.productionTip = false;

Vue.use(VueRouter);

const router = new VueRouter({
  routes,
  mode: 'history',
});

new Vue({
  router,
  render: h => h(App),
}).$mount('#app');

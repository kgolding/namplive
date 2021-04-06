import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';

Vue.config.productionTip = false;
Vue.config.devtools = true;

import * as Wails from '@wailsapp/runtime';

import { library } from '@fortawesome/fontawesome-svg-core'
import { faSortUp, faSortDown, faHeart, faEyeSlash, faSpinner } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(faSortUp, faSortDown, faHeart, faEyeSlash, faSpinner)

Vue.component('fa', FontAwesomeIcon)

Wails.Init(() => {
	new Vue({
		render: h => h(App)
	}).$mount('#app');
});

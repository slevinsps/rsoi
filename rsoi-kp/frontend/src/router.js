import Vue from 'vue';
import Router from 'vue-router';
import Home from './views/Home.vue';
import Login from './views/Login.vue';
import Register from './views/Register.vue';
import Monitor from './views/Monitor.vue';
import ChooseEquipments from './views/ChooseEquipments.vue';
import Equipment from './views/Equipment.vue';
import Model from './views/Model.vue';
import AdminAddEquipmentModel from './views/AdminAddEquipmentModel.vue'
import AdminAddEquipment from './views/AdminAddEquipment.vue'

Vue.use(Router);

export const router = new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/monitors',
      component: Home
    },
    {
      path: '/login',
      component: Login
    },
    {
      path: '/register',
      component: Register
    },
    {
      path: '/add_eqipment_model',
      name: 'add_eqipment_model',
      component: AdminAddEquipmentModel
    },
    {
      path: '/add_eqipment',
      name: 'add_eqipment',
      component: AdminAddEquipment
    },
    {
      path: '/equipment/:equipment_uuid',
      name:'equipment',
      component: Equipment,
      props: true,
    },
    {
      path: '/monitor/:monitor_uuid',
      name:'monitor',
      component: Monitor,
      props: true,
    },
    {
      path: '/model/:equipment_model_uuid',
      name:'model',
      component: Model,
      props: true,
    },
    {
      path: '/choose_equipments/:monitor_uuid',
      name:'choose_equipments',
      component: ChooseEquipments,
      props: true,
    },
    {
      path: '/admin',
      name: 'admin',
      // lazy-loaded
      component: () => import('./views/BoardAdmin.vue')
    }
  ]
});

// router.beforeEach((to, from, next) => {
//   const publicPages = ['/login', '/register', '/home'];
//   const authRequired = !publicPages.includes(to.path);
//   const loggedIn = localStorage.getItem('user');

//   // trying to access a restricted page + not logged in
//   // redirect to login page
//   if (authRequired && !loggedIn) {
//     next('/login');
//   } else {
//     next();
//   }
// });

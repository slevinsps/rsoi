<template>
  <div id="app">
    <nav class="navbar navbar-expand navbar-dark bg-dark">
      <a href class="navbar-brand" @click.prevent>Monitoring System</a>
      <div class="navbar-nav mr-auto">
        <li class="nav-item">
          <router-link to="/" class="nav-link">
            <font-awesome-icon icon="home" />Home
          </router-link>
        </li>
        <li v-if="userIsAdmin" class="nav-item">
          <router-link to="/admin" class="nav-link">Admin Board</router-link>
        </li>
      </div>

      
      <div v-if="!currentUser" class="navbar-nav ml-auto">
        <li class="nav-item">
          <router-link to="/register" class="nav-link">
            <font-awesome-icon icon="user-plus" />Sign Up
          </router-link>
        </li>
        <li class="nav-item">
          <router-link to="/login" class="nav-link">
            <font-awesome-icon icon="sign-in-alt" />Login
          </router-link>
        </li>
      </div>

      <div v-if="currentUser" class="navbar-nav ml-auto">
        <li class="nav-item">
          <a class="nav-link" href @click.prevent="logOut">
            <font-awesome-icon icon="sign-out-alt" />LogOut
          </a>
        </li>
      </div>
    </nav>

    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    },
    userIsAdmin() {
      if (this.currentUser) {
        console.log(this.currentUser)
        return this.currentUser.is_admin === 'true';
      }
      return false;
    },
    showModeratorBoard() {
      if (this.currentUser && this.currentUser.roles) {
        return this.currentUser.roles.includes('ROLE_MODERATOR');
      }

      return false;
    }
  },
  methods: {
    logOut() {
      this.$store.dispatch('auth/logout');
      this.$router.push('/login');
    },
  }
};
</script>


<style>
body {
  background-color: #f5f6f8;
}
.btn-primary {
  background-color: #343a40;
}
a {
  text-decoration: none;
}
</style>
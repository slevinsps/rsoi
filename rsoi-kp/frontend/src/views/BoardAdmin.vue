<template>
  <div class="container">
    
    <div class=title>Admin board </div>
    
    <router-link :to="{ name: 'add_eqipment'}" style="text-decoration:none;">      
      <div>Add equipment</div>
    </router-link>
    <router-link :to="{ name: 'add_eqipment_model'}" style="text-decoration:none;">      
      <div>Add equipment model</div>
    </router-link>
    
    
  </div>
</template>

<script>

export default {
  name: 'AdminBoard',
  data() {
    return {
      loading: false,
      message: ''
    };
  },
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
  },
  mounted() {
    if (!this.currentUser) {
      this.$router.push('/login');
      return
    }
    if (!this.userIsAdmin) {
      this.$router.push('/');
      return
    }
  },
};
</script>

<style scoped>
a {
  text-decoration: none;
}

.back {
  padding-top: 15px;
  padding-right: 15px;
  padding-left: 15px;
}

.back svg {
  font-size: 50;
}

.bullet {
margin-left: 0;
list-style: none;
counter-reset: li;
padding-inline-start: 0px
}
.bullet li {
position: relative;
margin-bottom: 1.5em;
border: 1px solid #343a40;
padding: 0.6em;
border-radius: 4px;
background: #fbfcff;
color: #231F20;
font-family: "Trebuchet MS", "Lucida Sans";
}
.bullet li:before {
position: absolute;
top: -0.7em;
padding-left: 0.4em;
padding-right: 0.4em;
font-size: 16px;
font-weight: bold;
color: #343a40;
background: #FEFEFE;
border-radius: 50%;
counter-increment: li;
content: counter(li);
}

.bullet li:hover {
  background: #ebeef5;
}

.button {
  padding: .5rem 0rem;
  color: #343a40
}

.button:hover {
  color: #007bff;
}

.btn-primary {
  background-color: #343a40;
  border-color: #343a40;
}

.title {
  padding: .5rem 0rem;
  font-size: 25px;
}
</style>
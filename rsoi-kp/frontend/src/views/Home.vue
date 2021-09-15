<template>
  <div class="container">
    <div class=title>Monitors</div>
    <a class="nav-link button" href @click.prevent="addMonitorMenu">
      <font-awesome-icon icon="plus" />Add monitor
    </a>
    <div class="form-group">
      <div v-if="message" class="alert alert-danger" role="alert">{{message}}</div>
    </div>
    <div class="add_monitor_menu">
      <form name="form" @submit.prevent="addMonitor">
        <div class="form-group">
          <label for="name">Name</label>
          <input
            v-model="monitor.name"
            v-validate="'required'"
            type="text"
            class="form-control"
            name="name"
          />
          <div
            v-if="errors.has('name')"
            class="alert alert-danger"
            role="alert"
          >Name is required!</div>
        </div>
        <div class="form-group">
          <button class="btn btn-primary btn-block " :disabled="loading">
            <span v-show="loading" class="spinner-border spinner-border-sm"></span>
            <span>Add</span>
          </button>
        </div>
      </form>
    </div>

    <ol class="bullet">
      <div v-if="monitors.length == 0" id="monitor_list">
        You dont have monitors
      </div>
      <div v-else id="monitor_list">
        <router-link v-for="(monitor, index) in monitors" :key="monitor.monitor_uuid" :to="{ name: 'monitor', params: { monitor_uuid: monitor.monitor_uuid, monitor_name: monitor.name } }" style="text-decoration:none;">
           <li>  
          <div class="box">
            <span id="left" class="equipmentInfo">{{ monitor.name }}</span>

            <span id="right" class="addButton" @click.prevent="deleteMonitor(monitor.monitor_uuid, index)"><font-awesome-icon icon="trash-alt" /></span>
            
           </div>
          </li>
        </router-link>
      </div>
    </ol>
  </div>
</template>

<script>
import MonitorService from '../services/monitor.service';
import Monitor from '../models/monitor';
// import ErrorHandler from '../services/errors.handler'

export default {
  name: 'Home',
  data() {
    return {
      monitor: new Monitor('', ''),
      monitors: [],
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
    }
  },
  mounted() {
    if (!this.currentUser) {
      this.$router.push('/login');
      return
    }

    let menu = document.getElementsByClassName('add_monitor_menu')[0];
    menu.style.visibility = 'hidden'
    MonitorService.getUserMonitors().then(
    response => {
        this.monitors = response.data;
        console.log(this.content)
      },
      error => {
        this.message =
                (error.response && error.response.data && error.response.data.message) ||
                error.message ||
                error.toString();
        console.log(this.message)
      }
    );
  },
  methods: {
    addMonitorMenu() {
      let menu = document.getElementsByClassName('add_monitor_menu')[0];
      let visibility = menu.style.visibility
      if (visibility == 'hidden') {
        menu.style.visibility = 'visible'
        menu.style.display = 'block'
      } else {
        menu.style.visibility = 'hidden'
        menu.style.display = 'none'
      }
    },
    addMonitor() {
      if (this.monitor.name) {
        MonitorService.addMonitor(this.monitor).then(
          response => {
              this.monitor = response.data;
              this.monitors.push(this.monitor);
            },
            error => {
              this.message =
                (error.response && error.response.data && error.response.data.message) ||
                error.message ||
                error.toString()
              console.log(error.toString())
            }
          );
      }
    }, 
    deleteMonitor(monitor_uuid, index) {
      if (monitor_uuid) {
        MonitorService.delMonitor(monitor_uuid).then(
          response => {
              this.$delete(this.monitors, index)
              console.log(index)
            }, error => {
              this.message =
                (error.response && error.response.data && error.response.data.message) ||
                error.message ||
                error.toString()
              console.log(error.toString())
            }
          );
      }
    }
  }
};
</script>

<style scoped>
a {
  text-decoration: none;
}


.box{
  height: 20px;
}
.addButton{
  height: 20px;
  width: 20px;
  text-align: center;
}
.addButton:hover{
  color: #2c66db;
}

#left{
     float:left;
}
#right{
     float:right;
}

.add_monitor_menu {
  display: none;
  visibility: hidden;
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
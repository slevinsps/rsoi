<template>
  <div class="container">
    <div class='row'>
      <div class='back' @click.prevent="back();"><font-awesome-icon icon="arrow-circle-left" /></div>
      <div class='title'>Equipment {{this.equipment.name}}</div> 
    </div>
    <div class="form-group">
      <div v-if="message" class="alert alert-danger" role="alert">{{message}}</div>
    </div>
    <div v-if="userIsAdmin" class="nav-item">
      <a class="nav-link " href @click.prevent="clearLogs">
        Clear logs for equipment
      </a>
      <div class="add_equipment_menu">
        <form name="form" @submit.prevent="changeStatus">
          <div>Status</div>
          <select id="selectStatus">
            <option disabled>Choose a status</option>
            <option value="ACTIVE">ACTIVE</option>
            <option value="INACTIVE">INACTIVE</option>
            <option value="DELETED">DELETED</option>
          </select>
          <div class="form-group">
            <button class="btn btn-primary btn-block " :disabled="loading">
              <span v-show="loading" class="spinner-border spinner-border-sm"></span>
              <span>Change status</span>
            </button>
          </div>
        </form>
      </div>
    </div>
    <div v-else class="nav-item">
     Status: {{this.equipment.status}}
    </div>


    <div class=description>
    <router-link :to="{ name: 'model', params: { equipment_model_uuid: this.equipment.equipment_model_uuid} }" style="text-decoration:none;">      
      <a class="nav-link button">
        Model: {{this.equipment.model_name}}
      </a>
    </router-link>
    </div>
    

    <ol class="bullet">
      <div id="monitor_list">
        <li v-for="record in data" :key="record.data_uuid" style="text-decoration:none;">
          temperature = {{ rounded(record.temperature) }}; voltage = {{ rounded(record.voltage) }}; frequency={{ rounded(record.frequency) }}; load level = {{ rounded(record.load_level) }}; time = {{ record.timestamp.split(".")[0] }}
        </li>
      </div>
    </ol>
  </div>
</template>

<script>
import EquipmentService from '../services/equipment.service';
import Equipment from '../models/equipment'

export default {
  name: 'Equipment',
  props: ['equipment_uuid'],
  data() {
    return {
      equipment: new Equipment('', '', '', '', ''),
      data: [],
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
    

    EquipmentService.getEquipment(this.equipment_uuid).then(
    response => {
        this.equipment = response.data;
        console.log(this.equipment)
        var e_status = document.getElementById("selectStatus");

        if (this.equipment.status == 'ACTIVE' || this.equipment.status == 'INACTIVE' || this.equipment.status == 'DELETED') {
          if (e_status) {
            e_status.value = this.equipment.status;
          }
        }
      },
      error => {
        this.message =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString()
        console.log(error.toString())
      }
    );
    

    EquipmentService.getData(this.equipment_uuid).then(
    response => {
        this.data = response.data;
        console.log(this.data)
      },
      error => {
        this.message =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString()
        console.log(error.toString())
      }
    );
  },
  methods: {
    back() {
      window.history.back();
    },
    rounded(number){
      return +number.toFixed(2);
    },
    clearLogs() {
      if (this.equipment.equipment_uuid) {
        EquipmentService.delData(this.equipment.equipment_uuid).then(
        response => {
            this.data = []      
            console.log(response.data)
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
    changeStatus() {
      var e_status = document.getElementById("selectStatus");
      this.equipment.status = e_status.value;
      EquipmentService.changeStatus(this.equipment).then(
      response => {
          if (response.data) {
            this.equipment.status = response.data.status
            var e_status = document.getElementById("selectStatus");
            if (this.equipment.status == 'ACTIVE' || this.equipment.status == 'INACTIVE' || this.equipment.status == 'DELETED') {
              e_status.value = this.equipment.status;
            }
          }
          console.log(response.data)
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
 
};
</script>



<style scoped>
.back {
  padding-top: 15px;
  padding-right: 15px;
  padding-left: 15px;
}

.back svg {
  font-size: 50;
}
a {
  text-decoration: none;
}

.description {
  padding: .5rem 0rem;
  font-size: 20px;
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
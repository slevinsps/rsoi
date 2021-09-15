<template>
  <div class="container">
    
    <div class='row'>
      <div class='back' @click.prevent="back();"><font-awesome-icon icon="arrow-circle-left" /></div>
      <div class=title>Monitor {{monitor_name}}</div>
    </div>

    <router-link :to="{ name: 'choose_equipments', params: { monitor_uuid: monitor_uuid} }" style="text-decoration:none;">      
      <a class="nav-link button">
        <font-awesome-icon icon="plus" />Choose equipment
      </a>
    </router-link>
    
    <ol class="bullet">
      <div id="equipment_list">
        <router-link v-for="(monitor_equipment, index) in monitor_equipments" :key="monitor_equipment.equipment_uuid" :to="{ name: 'equipment', params: { equipment_uuid: monitor_equipment.equipment_uuid } }" style="text-decoration:none;">
          <li>  
          <div class="box">
            <span id="left" class="equipmentInfo">Name: {{ monitor_equipment.name }}; Model: {{ monitor_equipment.model_name}}; Status: {{ monitor_equipment.status}}</span>

            <span id="right" class="addButton" @click.prevent="deleteEquipmentFromMonitor(monitor_equipment.equipment_uuid, index)"><font-awesome-icon icon="trash-alt" /></span>
            
           </div>
          </li>
        </router-link> 
      </div>
    </ol>
  </div>
</template>

<script>
import EquipmentService from '../services/equipment.service';
import Equipment from '../models/equipment';

export default {
  name: 'Monitor',
  props: ['monitor_uuid', 'monitor_name'],
  data() {
    return {
      monitor_equipment: new Equipment('', '', '', ''),
      monitor_equipments: [],
      loading: false,
      message: ''
    };
  },
  computed: {
    currentUser() {
      return this.$store.state.auth.user;
    }
  },
  mounted() {
    if (!this.currentUser) {
      this.$router.push('/login');
      return
    }
    EquipmentService.getMonitorEquipments(this.monitor_uuid).then(
    response => {
        this.monitor_equipments = response.data;
        console.log(this.monitor_equipments)
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
    getListOfEquipments() {
      if (this.monitor_equipment.name) {
        EquipmentService.addEquipmentToMonitor(this.monitor_equipment).then(
          response => {
              this.monitor_equipment = response.data;
              this.monitor_equipments.push(this.monitor_equipment);
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
    deleteEquipmentFromMonitor(equipment_uuid, index) {
      if (equipment_uuid && this.monitor_uuid) {
        EquipmentService.delEquipmentFromMonitor(equipment_uuid, this.monitor_uuid).then(
          response => {
              this.$delete(this.monitor_equipments, index)
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
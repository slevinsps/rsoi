<template>
  <div class="container">
    <div class='row'>
      <div class='back' @click.prevent="back();"><font-awesome-icon icon="arrow-circle-left" /></div>
      <div class='title'>Equipments</div> 
    </div>
    <div class='description'>Сhoose equipments for your monitor:</div>
    <div v-if="message" class="alert alert-danger" role="alert">{{message}}</div>
    <ol class="bullet">
      <div id="equipment_list">
        <li v-for="(equipment, index) in equipments" :key="equipment.equipment_uuid" style="text-decoration:none;">
           <div class="box">
            <span id="left" class="equipmentInfo">Name: {{ equipment.name }}; Model: {{ equipment.model_name}}; Status: {{ equipment.status}}</span>

            <span id="right" class="addButton" @click.prevent="addEquipmentToMonitor(equipment.equipment_uuid, index)"><font-awesome-icon icon="plus" /></span>
            
           </div>
        </li>
       
      </div>
    </ol>
  </div>
</template>

<script>
import EquipmentService from '../services/equipment.service';

export default {
  name: 'ChooseEquipments',
  props: ['monitor_uuid'],
  data() {
    return {
      equipments: [],
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
    EquipmentService.getNotAddedEquipments(this.monitor_uuid).then(
    response => {
        this.equipments = response.data;
        console.log(this.monitor_equipments)

      },
      error => {
        this.message =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString()
        console.log(this.message)
      }
    );
  },
  methods: {
    back() {
      window.history.back();
    },
    addEquipmentToMonitor(equipment_uuid, index) {
      if (equipment_uuid && this.monitor_uuid) {
        EquipmentService.addEquipmentToMonitor(equipment_uuid, this.monitor_uuid).then(
          response => {
              this.$delete(this.equipments, index)
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



.back {
  padding-top: 15px;
  padding-right: 15px;
  padding-left: 15px;
}

.back svg {
  font-size: 50;
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
.description {
  padding: .5rem 0rem;
  font-size: 20px;
}
</style>
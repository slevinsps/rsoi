<template>
  <div class="container">
    <div class=title>Add equipment model</div>
    <a class="nav-link button" href @click.prevent="showMenu">
      <font-awesome-icon icon="plus" />Add equipment model
    </a>
    <div class="form-group">
      <div v-if="message" class="alert alert-danger" role="alert">{{message}}</div>
    </div>
    <div class="add_equipment_model_menu">
      <form name="form" @submit.prevent="addEquipmentModel">
        <div class="form-group">
          <label for="name">Name</label>
          <input
            v-model="equipment_model.name"
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
          <input id="models_docs" name="myFile" type="file" accept=".pdf,.txt">
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
      <div v-if="equipment_models.length == 0" id="equipment_models_list">
        You dont have equipment models
      </div>
      <div v-else id="equipment_model_list">
        <!-- <li style="text-decoration:none;">
           {{ equipment_model.name }}
        </li> -->
        <router-link v-for="(equipment_model, index) in equipment_models" :key="equipment_model.equipment_model_uuid" :to="{ name: 'model', params: { equipment_model_uuid: equipment_model.equipment_model_uuid} }" style="text-decoration:none;">      
          <li>  
          <div class="box">
            <span id="left" class="equipmentInfo">Model name: {{ equipment_model.name }}</span>

            <span id="right" class="addButton" @click.prevent="deleteEquipmentModel(equipment_model.equipment_model_uuid, index)"><font-awesome-icon icon="trash-alt" /></span>
            
           </div>
          </li>
        </router-link>
      </div>
    </ol>
  </div>
</template>

<script>
import EquipmentService from '../services/equipment.service';
import EquipmentModel from '../models/equipment_model'
import DocumentService from '../services/documents.service';


export default {
  name: 'Home',
  data() {
    return {
      equipment_model: new EquipmentModel('', ''),
      equipment_models: [],
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
    if (!this.userIsAdmin) {
      this.$router.push('/');
      return
    }

    let menu = document.getElementsByClassName('add_equipment_model_menu')[0];
    menu.style.visibility = 'hidden'
    EquipmentService.getEquipmentModels().then(
    response => {
        this.equipment_models = response.data;
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
    showMenu() {
      let menu = document.getElementsByClassName('add_equipment_model_menu')[0];
      let visibility = menu.style.visibility
      if (visibility == 'hidden') {
        menu.style.visibility = 'visible'
        menu.style.display = 'block'
      } else {
        menu.style.visibility = 'hidden'
        menu.style.display = 'none'
      }
    },
    hasExtension(fileName, exts) {
      return (new RegExp('(' + exts.join('|').replace(/\./g, '\\.') + ')$')).test(fileName);
    },
    addEquipmentModel() {
      if (this.equipment_model.name) {
        EquipmentService.addEquipmentModel(this.equipment_model).then(
          response => {
              this.equipment_model = response.data;
              this.equipment_models.push(this.equipment_model);
              console.log(this.equipment_model.equipment_model_uuid)
              let files = document.getElementById("models_docs").files;
  
              if (this.equipment_model.equipment_model_uuid && files.length > 0) {
                let formData = new FormData();
                formData.append("file", files[0]);
                DocumentService.addFile(formData, this.equipment_model.equipment_model_uuid).then(
                  response => {
                      console.log('FIle was upload')
                    },
                    error => {
                      console.log('Error in DocumentService.addFile')
                      EquipmentService.delEquipmentModel(this.equipment_model.equipment_model_uuid).then(
                        response => {
                          console.log('Delete model')
                        },
                        error => {
                          this.message =
                            (error.response && error.response.data && error.response.data.message) ||
                            error.message ||
                            error.toString()
                          console.log(error.toString())
                        })
                      this.message =
                        (error.response && error.response.data && error.response.data.message) ||
                        error.message ||
                        error.toString()
                      console.log(error.toString())
                    }
                  );
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
      }
    },

    deleteEquipmentModel(equipment_model_uuid, index) {
      if (equipment_model_uuid) {
        EquipmentService.delEquipmentModel(equipment_model_uuid).then(
          response => {
              this.$delete(this.equipment_models, index)
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

.add_equipment_model_menu {
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
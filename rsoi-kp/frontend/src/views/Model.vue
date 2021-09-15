<template>
  <div class="container">
    
    <div class='row'>
      <div class='back' @click.prevent="back();"><font-awesome-icon icon="arrow-circle-left" /></div>
      <div class=title>Model {{this.equipment_model.name}}</div>
    </div>

    <ol class="bullet">
      <div v-if="documents.length == 0" id="document_list">
        No documents for this model
      </div>
      <div v-else id="document_list">
        <li v-for="document in documents" :key="document.file_uuid" :to="{ name: 'file', params: { file_uuid: document.file_uuid } }" style="text-decoration:none;" @click.prevent="downloadFile(document.file_uuid)">
          {{ document.name }}
        </li>
      </div>
    </ol>
  </div>
</template>

<script>
import DocumentService from '../services/documents.service';
import EquipmentService from '../services/equipment.service';

import EquipmentModel from '../models/equipment_model';

export default {
  name: 'Model',
  props: ['equipment_model_uuid'],
  data() {
    return {
      equipment_model: new EquipmentModel('', ''),
      documents: [],
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

    EquipmentService.getEquipmentModel(this.equipment_model_uuid).then(
    response => {
        this.equipment_model = response.data;
        console.log(this.documents)
      },
      error => {
        this.message =
          (error.response && error.response.data && error.response.data.message) ||
          error.message ||
          error.toString()
        console.log(error.toString())
      }
    );

    DocumentService.getDocuments(this.equipment_model_uuid).then(
    response => {
        this.documents = response.data;
        console.log(this.documents)
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
    downloadFile(file_uuid) {
      DocumentService.getFile(file_uuid).then(
        response => {
          const url = window.URL.createObjectURL(new Blob([response.data]));
          const link = document.createElement('a');
          link.href = url;
          link.setAttribute('download', 'file.pdf'); //or any other extension
          document.body.appendChild(link);
          link.click();
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
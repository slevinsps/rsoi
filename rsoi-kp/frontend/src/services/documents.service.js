import axios from 'axios';
import authHeader from './auth-header';


// const API_URL = 'http://localhost:8081/gateway/api/v1/app/';
const API_URL = '/gateway/api/v1/app/';


class DocumentService {
  getDocuments(equipmentModelUUID) {
    return axios.get(API_URL + 'documentation/equipment_model/' + equipmentModelUUID, { headers: authHeader() });
  }
  getFile(fileUUID) {
    return axios({
        url: API_URL + 'documentation/' + fileUUID, //your url
        method: 'GET',
        headers: authHeader(),
        responseType: 'blob', // important
    })
  }
  addFile(file, equipmentModelUUID) {
    // return axios.get(API_URL + 'documentation/' + fileUUID, { headers: authHeader() });
    let header = {}
    let user = JSON.parse(localStorage.getItem('user'));
    if (user && user.access_token) {
      header = { Authorization: 'Bearer ' + user.access_token,  'Content-Type': 'multipart/form-data'};
    }
    return axios.post(API_URL + 'documentation/upload/' + equipmentModelUUID, file, { headers: header})
  }
}


export default new DocumentService();

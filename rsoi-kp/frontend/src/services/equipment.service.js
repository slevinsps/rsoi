import axios from 'axios';
import authHeader from './auth-header';


// const API_URL = 'http://localhost:8081/gateway/api/v1/app/';
const API_URL = '/gateway/api/v1/app/';


class EquipmentService {
  getMonitorEquipments(monitor_uuid) {
    return axios.get(API_URL + 'equipment/list/' + monitor_uuid, { headers: authHeader() });
  }

  getNotAddedEquipments(monitor_uuid) {
    return axios.get(API_URL + 'equipment/list/' + monitor_uuid + '/notadded', { headers: authHeader() });
  }

  getEquipment(equipment_uuid) {
    return axios.get(API_URL + 'equipment/' + equipment_uuid, { headers: authHeader() });
  }

  getEquipmentModel(equipmentModelUUID) {
    return axios.get(API_URL + 'equipment/model/' + equipmentModelUUID, { headers: authHeader() });
  }

  getEquipmentModels() {
    return axios.get(API_URL + 'equipment/model/list', { headers: authHeader() });
  }

  getEquipments() {
    return axios.get(API_URL + 'equipment/list', { headers: authHeader() });
  }

  addEquipment(equipment) {
    return axios.post(API_URL + 'equipment/create', {name: equipment.name, equipment_model_uuid: equipment.equipment_model_uuid, status: equipment.status}, { headers: authHeader() });
  }

  changeStatus(equipment) {
    return axios.put(API_URL + 'equipment/update', {name: equipment.name, equipment_uuid: equipment.equipment_uuid, equipment_model_uuid: equipment.equipment_model_uuid, status: equipment.status}, { headers: authHeader() });
  }

  delEquipment(equipmentUUID) {
    return axios.delete(API_URL + 'equipment/' + equipmentUUID, { headers: authHeader() });
  }

  addEquipmentToMonitor(equipment_uuid, monitor_uuid) {
    return axios.post(API_URL + 'monitor/' + monitor_uuid + '/add/' + equipment_uuid, {}, { headers: authHeader() });
  }

  delEquipmentFromMonitor(equipment_uuid, monitor_uuid) {
    return axios.delete(API_URL + 'monitor/' + monitor_uuid + '/del/' + equipment_uuid,  { headers: authHeader() });
  }

  addEquipmentModel(equipment) {
    return axios.post(API_URL + 'equipment/model/create', {name: equipment.name}, { headers: authHeader() });
  }

  delEquipmentModel(equipmentModelUUID) {
    return axios.delete(API_URL + 'equipment/model/' + equipmentModelUUID, { headers: authHeader() });
  }

  getData(equipment_uuid) {
    return axios.get(API_URL + 'generator/equipment/' + equipment_uuid,  { headers: authHeader() });
  }

  delData(equipment_uuid) {
    return axios.delete(API_URL + 'generator/equipment/' + equipment_uuid,  { headers: authHeader() });
  }
}

export default new EquipmentService();

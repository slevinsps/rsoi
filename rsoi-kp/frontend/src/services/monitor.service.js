import axios from 'axios';
import authHeader from './auth-header';


// const API_URL = 'http://localhost:8081/gateway/api/v1/app/';
const API_URL = '/gateway/api/v1/app/';

class MonitorService {
  getUserMonitors() {
    return axios.get(API_URL + 'monitor/list', { headers: authHeader() });
  }

  addMonitor(monitor) {
    return axios.post(API_URL + 'monitor/create', {name: monitor.name}, { headers: authHeader() });
  }

  delMonitor(monitorUUID) {
    return axios.delete(API_URL + 'monitor/' + monitorUUID, { headers: authHeader() });
  }
}

export default new MonitorService();

export default class Monitor {
  constructor(data_uuid, equipment_uuid, temperature, voltage, frequency, load_level) {
    this.data_uuid = data_uuid;
    this.equipment_uuid = equipment_uuid;
    this.temperature = temperature;
    this.voltage = voltage;
    this.frequency = frequency;
    this.load_level = load_level;
  }
}

export default class Equipment {
  constructor(name, equipment_uuid, model_name, equipment_model_uuid, status) {
    this.name = name;
    this.equipment_uuid = equipment_uuid;
    this.model_name = model_name
    this.equipment_model_uuid = equipment_model_uuid;
    this.status = status;
  }
}

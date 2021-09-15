package api

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"services/config"
	"services/internal/models"
	"services/internal/repInterface"
	"services/utils"
	"time"

	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	rep         repInterface.RepInterface
	conf        *config.Configuration
	redisClient *redis.Client
	quit        chan struct{}
}

func NewHandler(rep repInterface.RepInterface, redisClient *redis.Client) *Handler {
	confPath := "conf.json"
	var (
		conf *config.Configuration
		err  error
	)
	if conf, err = config.Init(confPath); err != nil {
		return nil
	}
	quit := make(chan struct{})
	return &Handler{rep: rep, conf: conf, quit: quit, redisClient: redisClient}
}

func (h *Handler) Generate() (err error) {
	var (
		equipmentsArr []models.Equipment
		data          models.Data
		r_new         *http.Request
		code          int
		body          []byte
	)
	r_new, err = http.NewRequest("POST", h.conf.ServiceGateway.URL+"/check/user", bytes.NewBuffer(body))
	if code, body, err = h.SendRequest(r_new, h.conf.ServiceGateway.URL, "/equipment/list", "GET"); err != nil {
		utils.PrintDebug(err.Error())
		return
	}
	if code != http.StatusOK {
		utils.PrintDebug("Generate code != http.StatusOK code = ", code)
		return
	}

	if err = json.Unmarshal(body, &equipmentsArr); err != nil {
		utils.PrintDebug("Generate Error in json.Unmarshal")
		return
	}
	for i := 0; i < len(equipmentsArr); i++ {
		data.EquipmentUUID = equipmentsArr[i].EquipmentUUID
		data.Frequency = (rand.Float32() * 5) + 0.1
		data.LoadLevel = (rand.Float32() * 10)
		data.Temperature = (rand.Float32() * 150) + 20
		data.Voltage = (rand.Float32() * 200) + 20
		data.DataUUID = uuid.NewV4()
		data.Timestamp = time.Now()
		if err = h.rep.CreateData(data); err != nil {
			utils.PrintDebug(err.Error())
			continue
		}
	}
	return
}

func (h *Handler) GeneratorStart() {
	utils.PrintDebug("GeneratorStart()")
	close(h.quit)
	ticker := time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				h.Generate()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	h.quit = quit

	return
}

func (h *Handler) GeneratorStartHandler(rw http.ResponseWriter, r *http.Request) {
	utils.PrintDebug("GeneratorStartHandler()")
	h.GeneratorStart()
	return
}

func (h *Handler) GeneratorStopHandler(rw http.ResponseWriter, r *http.Request) {
	close(h.quit)
	h.quit = make(chan struct{})
	return
}

func (h *Handler) GetDataByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {

	const place = "GetDataByEquipmentUUID"

	var (
		err            error
		data_arr       []models.Data
		equipment_uuid uuid.UUID
	)

	rw.Header().Set("Content-Type", "application/json")

	if equipment_uuid, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if data_arr, err = h.rep.GetDataByEquipmentUUID(equipment_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetDataByEquipmentUUID")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(data_arr) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(data_arr)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) DeleteDataByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DeleteDataByEquipmentUUID"

	var (
		err            error
		equipment_uuid uuid.UUID
	)

	if equipment_uuid, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if err = h.rep.DeleteDataByEquipmentUUID(equipment_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	sendMessage(rw, "Data was deleted", http.StatusOK, nil)
	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) ClearAllData(rw http.ResponseWriter, r *http.Request) {
	const place = "ClearAllData"

	var (
		err error
	)

	if err = h.rep.CLearData(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	sendMessage(rw, "Data was cleared", http.StatusOK, nil)
	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) ServiceRegister(rw http.ResponseWriter, r *http.Request) {
	const place = "ServiceRegister"

	var (
		err     error
		service models.Service
		td      *models.TokenDetails
	)

	if service, err = getServiceFromBody(r); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		utils.PrintDebug("Error in getServiceFromBody(r) " + err.Error())
		return
	}
	if service.Login != os.Getenv("GATEWAY_ID") || service.Password != os.Getenv("GATEWAY_SECRET") {
		http.Error(rw, "you do not have permission to access this service", http.StatusUnauthorized)
		utils.PrintDebug("Error in checksecrets login = ", service.Login, "password = ", service.Password)
		return
	}

	td, err = CreateServiceToken(service.Login)
	if err != nil {
		utils.PrintDebug("Error in creating token " + err.Error())
		sendMessage(rw, "Error in creating token", http.StatusUnprocessableEntity, nil)
		return
	}
	rw.WriteHeader(http.StatusOK)
	resBytes, _ := json.Marshal(td)
	sendJSON(rw, resBytes)
	return
}

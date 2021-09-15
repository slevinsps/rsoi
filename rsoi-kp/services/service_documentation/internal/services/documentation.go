package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"services/config"
	"services/constants"
	"services/internal/models"
	"services/internal/repInterface"
	"services/utils"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	rep  repInterface.RepInterface
	conf *config.Configuration
}

func NewHandler(rep repInterface.RepInterface) *Handler {
	confPath := "conf.json"
	var (
		conf *config.Configuration
		err  error
	)
	if conf, err = config.Init(confPath); err != nil {
		return nil
	}
	return &Handler{rep: rep, conf: conf}
}

func ExtractToken(r *http.Request) string {
	fmt.Println("ExtractToken")
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	fmt.Println("ExtractToken")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	fmt.Println("VerifyToken")
	tokenString := ExtractToken(r)
	fmt.Println("doc get token ", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		utils.PrintDebug(err.Error())
		return nil, err
	}
	fmt.Println("VerifyToken")
	return token, nil
}

func TokenValid(r *http.Request) (token *jwt.Token, err error) {
	fmt.Println("TokenValid")
	token, err = VerifyToken(r)
	if err != nil {
		return
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return
	}
	fmt.Println("TokenValid")
	return
}

// TokenAuthMiddleware
func (h *Handler) TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)
		_, err = TokenValid(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			utils.PrintDebug("No valid token in TokenAuthMiddleware" + err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

func uploadFile(w http.ResponseWriter, r *http.Request) (error, string, string) {
	name := ""
	path := ""
	r.Body = http.MaxBytesReader(w, r.Body, constants.MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(constants.MAX_UPLOAD_SIZE); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 100MB in size", http.StatusBadRequest)
		return err, name, path
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err, name, path
	}

	defer file.Close()

	// Create the uploads folder if it doesn't
	// already exist
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err, name, path
	}

	// Create a new file in the uploads directory
	name = fileHeader.Filename
	path = fmt.Sprintf("./uploads/%d_%s", time.Now().UnixNano(), name)
	dst, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err, name, path
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err, name, path
	}

	return err, name, path
}

func sendFile(w http.ResponseWriter, r *http.Request, name string, path string) error {
	Openfile, err := os.Open(path)
	defer Openfile.Close()
	if err != nil {
		http.Error(w, "File not found.", 404)
		return err
	}

	FileHeader := make([]byte, 512)
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)
	return err
}

// CreateFile
func (h *Handler) CreateFile(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateFile"
	utils.PrintDebug(place)
	var (
		err                error
		file               models.File
		fileSend           models.FileSend
		name               string
		path               string
		equipmentModelUUID uuid.UUID
	)
	fmt.Println("AAAAAAAAAAAAAAAAaa")
	fileUID := uuid.NewV4()

	if err, name, path = uploadFile(rw, r); err != nil {
		return
	}

	if equipmentModelUUID, err = getUUID(r, "equipmentModelUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if len(name) == 0 {
		sendMessage(rw, "File name is empty", http.StatusBadRequest, nil)
		return
	}

	file.FileUUID = fileUID
	file.Name = name
	file.Path = path
	file.EquipmentModelUUID = equipmentModelUUID

	if err = h.rep.CreateFile(file); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateFile " + err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	fileSend.EquipmentModelUUID = file.EquipmentModelUUID
	fileSend.FileUUID = file.FileUUID
	fileSend.Name = file.Name
	rw.WriteHeader(http.StatusCreated)
	resBytes, _ := json.Marshal(fileSend)
	sendJSON(rw, resBytes)

	return
}

// UpdateFile
func (h *Handler) UpdateFile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateFile"
	utils.PrintDebug(place)
	var (
		err           error
		file_new      models.File
		file_old      models.File
		file_res      models.File
		checkFindFile bool
	)

	if file_new, err = getFileFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getFileFromBody " + err.Error())
		return
	}
	file_uuid := file_new.FileUUID

	if file_old, checkFindFile, err = h.rep.GetFileByFileUUID(file_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if !checkFindFile {
		sendMessage(rw, "Can't find file by uuid: "+file_uuid.String(), http.StatusNotFound, nil)
	}

	if file_new.Name == "" {
		file_new.Name = file_old.Name
	}
	if file_new.EquipmentModelUUID == uuid.Nil {
		file_new.EquipmentModelUUID = file_old.EquipmentModelUUID
	}

	if file_res, err = h.rep.UpdateFile(file_new); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateFile " + err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusCreated)
	resBytes, _ := json.Marshal(file_res)
	sendJSON(rw, resBytes)

	return
}

func (h *Handler) GetAllFiles(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllUserFiles"

	var (
		err   error
		files []models.File
	)

	rw.Header().Set("Content-Type", "application/json")

	if files, err = h.rep.GetAllFiles(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetAllFiles")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(files) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(files)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) GetAllFilesByEquipmentModelUUID(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllFilesByEquipmentModelUUID"

	var (
		err                error
		files              []models.File
		equipmentModelUUID uuid.UUID
	)

	rw.Header().Set("Content-Type", "application/json")

	if equipmentModelUUID, err = getUUID(r, "equipmentModelUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if files, err = h.rep.GetAllFilesByEquipmentModelUUID(equipmentModelUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetAllFiles")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(files) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(files)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) GetFileByFileUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetFileByUUID"

	var (
		err           error
		file_uuid     uuid.UUID
		file          models.File
		checkFindFile bool
	)

	rw.Header().Set("Content-Type", "application/json")

	if file_uuid, err = getUUID(r, "fileUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if file, checkFindFile, err = h.rep.GetFileByFileUUID(file_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if checkFindFile {
		if err = sendFile(rw, r, file.Name, file.Path); err != nil {
			return
		}
	} else {
		sendMessage(rw, "Can't find file by uuid: "+file_uuid.String(), http.StatusNotFound, nil)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) DelFileByFileUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelFileByFileUUID"

	var (
		err           error
		file_uuid     uuid.UUID
		checkFindFile bool
	)

	if file_uuid, err = getUUID(r, "fileUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if _, checkFindFile, err = h.rep.GetFileByFileUUID(file_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindFile {
		sendMessage(rw, "Can't find file by uuid: "+file_uuid.String(), http.StatusNotFound, nil)
		return
	}

	if err = h.rep.DelFileByFileUUID(file_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	sendMessage(rw, "File was deleted", http.StatusOK, nil)
	printResult(err, http.StatusOK, place)
	return
}

func CreateServiceToken(login string) (*models.TokenDetails, error) {
	var err error

	td := &models.TokenDetails{}

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["source"] = "service"
	atClaims["service login"] = login
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
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
		utils.PrintDebug("Error in checksecrets")
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

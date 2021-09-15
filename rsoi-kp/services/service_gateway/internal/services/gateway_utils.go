package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"services/utils"
	"strings"
	"time"

	"services/internal/models"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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
	fmt.Println(tokenString)
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

func (h *Handler) ExtractTokenMetadata(r *http.Request) (ad models.AccessDetails, err error) {
	fmt.Println("ExtractTokenMetadata")
	var (
		userId uuid.UUID
	)

	token, err := VerifyToken(r)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return
		}
		userId, err = uuid.FromString(fmt.Sprintf("%s", claims["user_uuid"]))
		if err != nil {
			return
		}
		ad.AccessUuid = accessUuid
		ad.UserId = userId
		return ad, nil
	}
	fmt.Println("ExtractTokenMetadata")
	return
}

func (h *Handler) TokenAuthMiddleware(next http.Handler, adminOnly bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			code     int
			body_new []byte
			td       models.TokenDetails
			tokenJWT *jwt.Token
			err      error
			source   string
			ok       bool
			isAdmin  bool
		)
		tokenJWT, err = TokenValid(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			utils.PrintDebug("No valid token in TokenAuthMiddleware" + err.Error())
			return
		}
		claims, ok := tokenJWT.Claims.(jwt.MapClaims)
		if ok && tokenJWT.Valid {
			source, ok = claims["source"].(string)
		}

		token := ExtractToken(r)
		td.AccessToken = token

		if body_new, err = json.Marshal(td); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			utils.PrintDebug("Error in Marshal(td) " + err.Error())
			return
		}
		if source == "user" {
			isAdmin, ok = claims["isAdmin"].(bool)
			utils.PrintDebug("adminOnly = ", adminOnly)
			utils.PrintDebug("isAdmin = ", isAdmin)
			if adminOnly && !isAdmin {
				http.Error(w, "unauthorized you are not admin", http.StatusForbidden)
				return
			}

			r_new, err := http.NewRequest("POST", h.conf.ServiceSession.URL+"/check/user", bytes.NewBuffer(body_new))
			if code, _, err = h.SendRequest(r_new, h.conf.ServiceSession.URL, "/check/user", "GET"); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				utils.PrintDebug("Error in SendRequest " + err.Error())
				return
			}

			if code != http.StatusOK {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) RegisterService(serviceBase string) (token string, err error) {
	var (
		body_new []byte
		service  models.Service
		td       models.TokenDetails
		resp     *http.Response
	)

	service.Login = os.Getenv("GATEWAY_ID")
	service.Password = os.Getenv("GATEWAY_SECRET")
	body_new, err = json.Marshal(service)
	if err != nil {
		utils.PrintDebug("Error in Marshal(service) " + err.Error())
		return
	}

	r_new := new(http.Request)
	r_new.Body = ioutil.NopCloser(bytes.NewBuffer(body_new))

	fmt.Println(serviceBase + "/service/register")
	if resp, err = utils.SendRequest(serviceBase+"/service/register", r_new, "POST"); err != nil {
		utils.PrintDebug("Error in SendRequest(/service/register " + err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		utils.PrintDebug("Error in registration service")
		err = errors.New("Error in registration service")
		return
	}
	if body_new, err = ioutil.ReadAll(resp.Body); err != nil {
		utils.PrintDebug("Error in ioutil.ReadAll " + err.Error())
		return
	}
	if err = json.Unmarshal(body_new, &td); err != nil {
		utils.PrintDebug("Error in json.Unmarshal(body_new, &td) " + err.Error())
		return
	}

	if err = h.redisClient.Set(serviceBase, td.AccessToken, 0).Err(); err != nil {
		utils.PrintDebug("Error in h.redisClient.Set " + err.Error())
		return
	}
	token = td.AccessToken
	return
}

func (h *Handler) SendRequest(r *http.Request, serviceBase string, urlString string, method string) (code int, body []byte, err error) {
	var (
		resp  *http.Response
		token string
	)

	if token, err = h.redisClient.Get(serviceBase).Result(); err != nil {
		utils.PrintDebug("Error in h.redisClient.Get " + err.Error())
		if token, err = h.RegisterService(serviceBase); err != nil {
			utils.PrintDebug("h.RegisterService " + err.Error())
			code = http.StatusInternalServerError
			return
		}
	}
	fmt.Println(token)
	r.Header.Set("Authorization", "B "+token)

	fmt.Println("a ", serviceBase+urlString)
	if resp, err = utils.SendRequest(serviceBase+urlString, r, method); err != nil {
		return
	}
	fmt.Println("b ", serviceBase+urlString)
	defer resp.Body.Close()
	code = resp.StatusCode
	fmt.Println("code = ", code)

	if code == http.StatusUnauthorized {
		fmt.Println("Registration service in gateway ")
		if token, err = h.RegisterService(serviceBase); err != nil {
			code = http.StatusInternalServerError
			return
		}
		r.Header.Set("Authorization", "B "+token)
		if resp, err = utils.SendRequest(serviceBase+urlString, r, method); err != nil {
			if resp.StatusCode == http.StatusUnauthorized {
				code = http.StatusInternalServerError
			}
			return
		}
	}
	fmt.Println("FFFFF")
	body, err = ioutil.ReadAll(resp.Body)

	return
}

func (h *Handler) CheckUser(rw http.ResponseWriter, r *http.Request) (userUUID uuid.UUID, err error) {
	var (
		ad       models.AccessDetails
		code     int
		body     []byte
		body_new []byte
		td       models.TokenDetails
	)

	token := ExtractToken(r)
	td.AccessToken = token
	if body_new, err = json.Marshal(td); err != nil {
		http.Error(rw, "unauthorized", http.StatusUnauthorized)
		utils.PrintDebug("Error in Marshal(td) " + err.Error())
		return
	}
	r_new := new(http.Request)
	r_new.Body = ioutil.NopCloser(bytes.NewBuffer(body_new))
	r_new.Header = r.Header

	if code, body, err = h.SendRequest(r_new, h.conf.ServiceSession.URL, "/check/user", "GET"); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		utils.PrintDebug("Error in SendRequest " + err.Error())
		return
	}
	fmt.Println(code)

	if code != http.StatusOK {
		rw.WriteHeader(code)
		sendJSON(rw, body)
		err = errors.New("no such user")
		return
	}

	if err = json.Unmarshal(body, &ad); err != nil {
		return
	}
	userUUID = ad.UserId
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
		err              error
		service          models.Service
		serviceFound     models.Service
		td               *models.TokenDetails
		checkFindService bool
	)

	if service, err = getServiceFromBody(r); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		utils.PrintDebug("Error in getServiceFromBody(r) " + err.Error())
		return
	}
	if serviceFound, checkFindService, err = h.db.GetServiceByLogin(service.Login); err != nil {
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		utils.PrintDebug("Error in GetServiceByLogin login = ", service.Login, "password = ", service.Password)
		return
	}

	if !checkFindService {
		http.Error(rw, "you do not have permission to access this service", http.StatusUnauthorized)
		utils.PrintDebug("No found such user = ", service.Login)
	}

	if !CheckPasswordHash(service.Password, serviceFound.Password) {
		http.Error(rw, "you do not have permission to access this service", http.StatusUnauthorized)
		utils.PrintDebug("NO CheckPasswordHash login = ", service.Login, "password = ", service.Password)
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

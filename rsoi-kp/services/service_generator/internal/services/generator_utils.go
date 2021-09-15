package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"services/internal/models"
	"services/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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
	fmt.Println(td.AccessToken)
	if err != nil {
		return nil, err
	}

	return td, nil
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

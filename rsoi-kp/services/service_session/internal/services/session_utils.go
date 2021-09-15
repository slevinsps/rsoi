package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"services/internal/models"
	"services/utils"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func CreateToken(userUUID uuid.UUID, isAdmin bool) (*models.TokenDetails, error) {
	var err error

	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 60).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["source"] = "user"
	atClaims["isAdmin"] = isAdmin
	atClaims["user_uuid"] = userUUID
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_uuid"] = userUUID
	rtClaims["isAdmin"] = isAdmin
	rtClaims["source"] = "user"
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractToken(r *http.Request, mode string) string {
	if mode == "header" {
		bearToken := r.Header.Get("Authorization")
		strArr := strings.Split(bearToken, " ")
		if len(strArr) == 2 {
			return strArr[1]
		}
	} else if mode == "body" {
		var (
			td  models.TokenDetails
			err error
		)
		if td, err = getTDFromBody(r); err != nil {
			utils.PrintDebug("Error in getTDFromBody(r)")
			return ""
		}
		return td.AccessToken
	}
	return ""
}

func VerifyToken(r *http.Request, mode string) (*jwt.Token, error) {
	tokenString := ExtractToken(r, mode)
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			utils.PrintDebug("Error in token.Method.(*jwt.SigningMethodHMAC)")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		utils.PrintDebug("Error in jwt.Parse(tokenString")
		return nil, err
	}
	return token, nil
}

func TokenValid(r *http.Request, mode string) error {
	token, err := VerifyToken(r, mode)
	if err != nil {
		utils.PrintDebug("Error in VerifyToken(r, mode)")
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		utils.PrintDebug("Error in token.Claims.(jwt.Claims)")
		return err
	}
	return nil
}

func ExtractTokenMetadata(r *http.Request, mode string) (*models.AccessDetails, error) {
	token, err := VerifyToken(r, mode)
	if err != nil {
		utils.PrintDebug("Error in VerifyToken(r, mode)")
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			utils.PrintDebug("Error in claims[access_uuid].(string)")
			return nil, err
		}
		userId, err := uuid.FromString(fmt.Sprintf("%s", claims["user_uuid"]))
		if err != nil {
			utils.PrintDebug("Error in uuid.FromString")
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func (h *Handler) RefreshTokens(rw http.ResponseWriter, r *http.Request) {
	var (
		td  models.TokenDetails
		err error
	)
	if td, err = getTokenDetailsFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getTokenDetatilsFromBody " + err.Error())
		return
	}
	refreshToken := td.RefreshToken

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		sendMessage(rw, "Refresh token expired", http.StatusUnauthorized, nil)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		sendMessage(rw, "Refresh token invalid", http.StatusUnauthorized, nil)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			sendMessage(rw, "Error in Claim refresh token", http.StatusUnprocessableEntity, nil)
			return
		}

		userId, err := uuid.FromString(fmt.Sprintf("%s", claims["user_uuid"]))
		if err != nil {
			sendMessage(rw, "Error in uuid refresh token", http.StatusUnprocessableEntity, nil)
			return
		}

		if claims["isAdmin"] == nil {
			claims["isAdmin"] = false
		}
		isAdmin := claims["isAdmin"].(bool)

		deleted, delErr := h.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 {
			sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
			return
		}

		ts, createErr := CreateToken(userId, isAdmin)
		if createErr != nil {
			sendMessage(rw, createErr.Error(), http.StatusForbidden, nil)
			return
		}
		saveErr := h.CreateAuth(userId, ts)
		if saveErr != nil {
			sendMessage(rw, createErr.Error(), http.StatusForbidden, nil)
			return
		}

		rw.WriteHeader(http.StatusCreated)
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		resBytes, _ := json.Marshal(tokens)
		sendJSON(rw, resBytes)
	} else {
		sendMessage(rw, "refresh expired", http.StatusUnauthorized, nil)
	}
}

func (h *Handler) FetchAuth(authD *models.AccessDetails) (userID uuid.UUID, err error) {
	var (
		useridStr string
	)
	useridStr, err = h.redisClient.Get(authD.AccessUuid).Result()
	if err != nil {
		utils.PrintDebug("Error in h.redisClient.Get(authD.AccessUuid)")
		return
	}
	userID, err = uuid.FromString(useridStr)
	if err != nil {
		utils.PrintDebug("Error in uuid.FromString(useridStr)")
		return
	}
	return
}

func (h *Handler) CreateAuth(userid uuid.UUID, td *models.TokenDetails) error {
	stringUserID := userid.String()
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := h.redisClient.Set(td.AccessUuid, stringUserID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := h.redisClient.Set(td.RefreshUuid, stringUserID, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (h *Handler) DeleteAuth(givenUuid string) (deleted int64, err error) {
	deleted, err = h.redisClient.Del(givenUuid).Result()
	if err != nil {
		return
	}
	return
}

func (h *Handler) TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := TokenValid(r, "header")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

package coinmonserver

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/VacantionIT/coin-price-mon/internal/app/store"
	"github.com/dgrijalva/jwt-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

// CoinMonServer - структура для сервера
type CoinMonServer struct {
	config *Config
	router *router.Router
	store  *store.Store
}

// New - Метод создание и инициализации нового сервера
func New(config *Config) *CoinMonServer {
	return &CoinMonServer{
		config: config,
		router: router.New(),
	}
}

// Start - Метод запуска сервера
func (s *CoinMonServer) Start() error {

	log.Print("coin server started...")

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	return fasthttp.ListenAndServe(s.config.BindAddr, s.router.Handler)
}

func (s *CoinMonServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st
	return nil
}

// SignIn - метод создания токена
func (s *CoinMonServer) SignIn(userName, userPass string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(userPass))
	pwd.Write([]byte(s.config.HashSalt))
	// password_hash := fmt.Sprintf("%x", pwd.Sum(nil))

	var myToken *jwt.Token

	// # через окружение
	testUser := os.Getenv("COIN_SERVER_USER_NAME")
	testPass := os.Getenv("COIN_SERVER_PASSWORD")
	if userName == testUser && userPass == testPass {
		fmt.Print("ok")
		t_now := jwt.TimeFunc().Unix()
		t_exp := jwt.TimeFunc().Add(time.Second * time.Duration(s.config.ExpDuration)).Unix()
		myToken = jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: t_exp,
				IssuedAt:  t_now,
			},
			Username: userName,
		})
	} else {
		return "", ErrUserDoesNotExist
	}
	return myToken.SignedString([]byte(s.config.TokenKey))
}

// CheckToken - метод проверки токена на основе данных в реквесте
func (s *CoinMonServer) CheckToken(ctx *fasthttp.RequestCtx) (bool, error) {
	// fmt.Print(ctx.Request.Header)
	authHeader := string(ctx.Request.Header.Peek("Authorization"))
	if len(authHeader) == 0 {
		return false, ErrInvalidAccessToken
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return false, ErrInvalidAccessToken
	}

	if headerParts[0] != "Bearer" {
		return false, ErrInvalidAccessToken
	}

	_, err := s.ParseToken(headerParts[1])
	// fmt.Print(user)
	if err != nil {
		log.Print(err)
		return false, ErrInvalidAccessToken
	}
	return true, nil
}

// ParseToken - Метода проверки токена. Возвращает имя пользователя из токена.
func (s *CoinMonServer) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.TokenKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.Username, nil
	}

	return "", ErrInvalidAccessToken
}

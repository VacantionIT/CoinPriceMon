package coinmonserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type response struct {
	Status string `json:"status"`
	Msg    string `json:"message,omitempty"`
	Token  string `json:"token,omitempty"`
}

func sendResponse(ctx *fasthttp.RequestCtx, status, msg, token string) {
	answer := &response{
		Status: status,
		Msg:    msg,
		Token:  token,
	}
	b, err := json.Marshal(answer)
	if err != nil {
		log.Printf("Error: %s", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		fmt.Fprint(ctx, string(b))
	}

}

func (s *CoinMonServer) handlerHello(ctx *fasthttp.RequestCtx) {
	sendResponse(ctx, STATUS_OK, "hello", "")

}

func (s *CoinMonServer) handlerAddCoin(ctx *fasthttp.RequestCtx) {
	isValid, err := s.CheckToken(ctx)
	if !isValid {
		if err == ErrInvalidAccessToken {
			sendResponse(ctx, STATUS_ERROR, "Unauthorized", "")
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		} else {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		}
	} else {
		coinID := string(ctx.FormValue("coin_id"))
		if coinID == "" {
			sendResponse(ctx, STATUS_ERROR, "coin_id is required", "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		interval, err := strconv.Atoi(string(ctx.FormValue("interval")))
		if err != nil {
			sendResponse(ctx, STATUS_ERROR, "bad interval, must be integer (in sec)", "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		if interval < 1 {
			sendResponse(ctx, STATUS_ERROR, "bad interval, must be greater then 1", "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		collection := s.store.DBClient.Database("coinmonserver").Collection("monitoring")
		filter := bson.D{primitive.E{Key: "coinid", Value: coinID}}

		update := bson.D{primitive.E{
			Key: "$set", Value: bson.D{primitive.E{
				Key: "interval", Value: interval},
			}},
		}
		_, err = collection.UpdateOne(
			// dbCtx,
			context.TODO(),
			filter,
			update,
			options.Update().SetUpsert(true),
		)
		if err != nil {
			log.Print(err)
			sendResponse(ctx, STATUS_ERROR, "db insert/update error", "")
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			return
		}
		sendResponse(ctx, STATUS_OK, coinID, "")
	}
}

func (s *CoinMonServer) handlerSignIn(ctx *fasthttp.RequestCtx) {
	name := string(ctx.FormValue("username"))
	passw := string(ctx.FormValue("password"))
	if name == "" || passw == "" {
		sendResponse(ctx, STATUS_ERROR, "bad parameters", "")
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	token, err := s.SignIn(name, passw)
	if err != nil {
		log.Print(name, err)

		if err == ErrUserDoesNotExist {
			sendResponse(ctx, STATUS_ERROR, err.Error(), "")
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}
		sendResponse(ctx, STATUS_ERROR, "we are working on this error", "")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	sendResponse(ctx, STATUS_OK, "", token)
}

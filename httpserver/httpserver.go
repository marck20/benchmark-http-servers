package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

// Each user struct emulates a database entry
type user struct {
	ID   int    `json:"ID"`
	Name string `json:"name,omitempty"`
	Role string `json:"role,omitempty"`
}

var (
	addr    = flag.String("addr", ":8080", "TCP address to listen to")
	usersDB = map[int]user{1: user{1, "admin", "admin"}, 2: user{2, "marc", "user"}, 3: user{3, "jordi", "user"}} //Emulates a database with three fields: id (int), name (string), role (string)
)

// Recives HTTP requests that target /user path and returns
// its corresponding user entry on usersDB.
// This endpoint expects an id (integer) argument which
// represents the user id.
// @param ctx *fasthttp.RequestCtx
func userHandler(ctx *fasthttp.RequestCtx) {
	userID, err := ctx.QueryArgs().GetUint("id") //Get id argument
	if err != nil {
		json.NewEncoder(ctx).Encode(map[string]string{"Error": "undefined id parameter."})
	} else {
		if user, ok := usersDB[userID]; ok {
			json.NewEncoder(ctx).Encode(user) //Encode requested user entry in a JSON object
		} else {
			json.NewEncoder(ctx).Encode(map[string]string{"Error": "Invalid user id."})
		}
	}

	ctx.SetContentType("application/json; charset=utf8")
}

//Simple endpoint which responds with a greeting message
// @param ctx *fasthttp.RequestCtx
func greetHandler(ctx *fasthttp.RequestCtx) {
	json.NewEncoder(ctx).Encode(map[string]string{"Message": "Hello!"})
	ctx.SetContentType("application/json; charset=utf8")
}

// Handles requests to /time and returns current time in a JSON.
// @param ctx *fasthttp.RequestCtx
func timeHandler(ctx *fasthttp.RequestCtx) {
	dt := time.Now()
	json.NewEncoder(ctx).Encode(map[string]string{"time": dt.String()}) //Encode current time in a JSON object
	ctx.SetContentType("application/json; charset=utf8")
}

// Main handler recive all requests and forwards them
// to its corresponding handler. Offers three endpoints:
// /user, /greet and /time.
// @param ctx *fasthttp.RequestCtx
func requestHandler(ctx *fasthttp.RequestCtx) {
	connectionClose := ctx.Request.ConnectionClose()
	if connectionClose == true {
		ctx.SetConnectionClose()
	}

	switch string(ctx.Path()) {
	case "/user":
		userHandler(ctx)
	case "/greet":
		greetHandler(ctx)
	case "/time":
		timeHandler(ctx)
	}
}

func main() {
	flag.Parse()

	h := requestHandler
	if err := fasthttp.ListenAndServe(*addr, h); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

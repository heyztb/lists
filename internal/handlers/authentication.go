package handlers

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/1Password/srp"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/jwt"
	"github.com/heyztb/lists-backend/internal/middleware"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type authenticationResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type authenticationSocketMessage struct {
	// current step of auth handshake
	Serial int `json:"serial"`
	// identifier
	I string `json:"I,omitempty"`
	// big A, client's ephemeral public key
	A *big.Int `json:"A,omitempty"`
	// client M, their proof of the generated session key
	C_M *big.Int `json:"C_M,omitempty"`
}

type authenticationSocketResponse struct {
	// current step of auth handshake
	Serial int `json:"serial"`
	// user salt
	S string `json:"s,omitempty"`
	// big B, server's ephemeral public key
	B     *big.Int `json:"B,omitempty"`
	Token string   `json:"user_id,omitempty"`
}

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	var (
		I               string
		s, v, A, B, C_M *big.Int
		k               []byte
		srpServer       *srp.SRP
		user            *models.User
	)
	_, ok := r.Context().Value(middleware.SessionKeyCtxKey).(string)
	if ok {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, &authenticationResponse{
			Status:  http.StatusOK,
			Message: "Authenticated",
		})
		return
	}
	upgrader := websocket.Upgrader{
		HandshakeTimeout: 30,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			render.Status(r, status)
			render.JSON(w, r, &authenticationResponse{
				Status:  status,
				Message: reason.Error(),
			})
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("websocket upgrader failure")
		return
	}
	ws.WriteJSON(&authenticationSocketResponse{
		Serial: 1,
	})
messageLoop:
	for {
		msg := &authenticationSocketMessage{}
		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Err(err).Msg("error reading incoming websocket message")
			break messageLoop
		}
		switch msg.Serial {
		// handle closes from the client
		case -1:
			break messageLoop
		case 1:
			I, A = msg.I, msg.A
			user, err = models.Users(qm.Where("username = ?", I), qm.Load(models.UserRels.Setting)).One(r.Context(), database.DB)
			if err != nil {
				log.Err(err).Msg("failed to fetch user from database")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			s, v = srp.NumberFromString(user.Salt), srp.NumberFromString(user.Verifier)
			// Choosing to use the 3072 group because of this formula that estimates actual key strength of k -- 8 * (SRP.Group.ExponentSize / 2)
			// in the case of 3072, this works out to 8 * (32/2) = 128 bits. keys are always 256 bits in length
			srpServer = srp.NewServerStd(srp.KnownGroups[srp.RFC5054Group3072], v)
			if srpServer == nil {
				log.Error().Msg("failed to initialize srp server component")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			err = srpServer.SetOthersPublic(A)
			if err != nil {
				log.Err(err).Msg("failed to set other parties ephemeral public key")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			B = srpServer.EphemeralPublic()
			ws.WriteJSON(&authenticationSocketResponse{
				Serial: msg.Serial + 1,
				S:      user.Salt,
				B:      B,
			})
		case 2:
			k, err = srpServer.Key()
			if err != nil {
				log.Err(err).Msg("failed to generate session key")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			ws.WriteJSON(&authenticationSocketResponse{
				Serial: msg.Serial + 1,
			})
		case 3:
			C_M = msg.C_M
			// compare the server proof with the one we receive the from the client
			if !srpServer.GoodServerProof(s.Bytes(), I, C_M.Bytes()) {
				log.Error().Msg("received bad proof from client")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			settings := user.R.Setting
			// cache the session key in redis -- ideally we will refresh the expiration timer on use, so that as long as the user is active we can continue
			// to send and receive messages as normal, and only after the redis entry truly expires will they need to reauth
			err := database.Redis.SetEx(
				r.Context(),
				fmt.Sprintf("%s:%d", database.RedisSessionKeyPrefix, user.ID),
				hex.EncodeToString(k),
				time.Duration(settings.SessionDuration)*time.Second,
			).Err()
			if err != nil {
				log.Err(err).Msg("failed to store session key in redis")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			token, err := jwt.GenerateToken(user.ID, uint64(settings.SessionDuration))
			if err != nil {
				log.Err(err).Uint64("id", user.ID).Msg("failed to generate jwt token for user")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				break messageLoop
			}
			ws.WriteJSON(&authenticationSocketResponse{
				Serial: 9999,
				Token:  token,
			})
			break messageLoop
		}
	}
	ws.Close()
}

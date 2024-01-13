package handlers

import (
	"encoding/hex"
	"math/big"
	"net/http"
	"time"

	"github.com/1Password/srp"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"github.com/heyztb/lists-backend/internal/database"
	"github.com/heyztb/lists-backend/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type authenticationResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type authenticationSocketMessage struct {
	Serial int      `json:"serial"`
	I      string   `json:"I,omitempty"`
	A      *big.Int `json:"A,omitempty"`
	C_M    *big.Int `json:"C_M,omitempty"`
}

type authenticationSocketResponse struct {
	Serial int      `json:"serial"`
	S      string   `json:"s,omitempty"`
	B      *big.Int `json:"B,omitempty"`
}

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	var (
		userID          int
		I               string
		s, v, A, B, C_M *big.Int
		k               []byte
	)
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
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, &authenticationResponse{
			Status:  http.StatusBadRequest,
			Message: "Failed to create websocket connection",
		})
	}
	ws.WriteJSON(&authenticationSocketResponse{
		Serial: 1,
	})
	for {
		msg := &authenticationSocketMessage{}
		err = ws.ReadJSON(&msg)
		if err != nil {
			log.Err(err).Msg("error reading incoming websocket message")
			return
		}

		switch msg.Serial {
		case -1:
			ws.Close()
			return
		case 1:
			I, A = msg.I, msg.A
			user, err := models.Users(qm.Where("username = ?", I)).One(r.Context(), database.DB)
			if err != nil {
				log.Err(err).Msg("failed to fetch user from database")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				return
			}
			userID = user.UserID
			s, v = srp.NumberFromString(user.Salt), srp.NumberFromString(user.Verifier)
			srpServer := srp.NewServerStd(srp.KnownGroups[srp.RFC5054Group3072], v)
			err = srpServer.SetOthersPublic(A)
			if err != nil {
				log.Err(err).Msg("failed to set other parties ephemeral public key")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				return
			}
			B = srpServer.EphemeralPublic()
			ws.WriteJSON(&authenticationSocketResponse{
				Serial: msg.Serial + 1,
				S:      user.Salt,
				B:      B,
			})
		case 2:
			srpServer := srp.NewServerStd(srp.KnownGroups[srp.RFC5054Group3072], v)
			// not checking for error because this will be the 2nd time during the handshake that we handle this A value, if we get to this point we know that A is valid
			srpServer.SetOthersPublic(A)
			k, err = srpServer.Key()
			if err != nil {
				log.Err(err).Msg("failed to generate session key")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				return
			}
			ws.WriteJSON(&authenticationSocketResponse{
				Serial: msg.Serial + 1,
			})
		case 3:
			C_M = msg.C_M
			srpServer := srp.NewServerStd(srp.KnownGroups[srp.RFC5054Group3072], v)
			srpServer.SetOthersPublic(A)
			srpServer.EphemeralPublic().Set(B)
			// compare the server proof with the one we receive the from the client
			if !srpServer.GoodServerProof(s.Bytes(), I, C_M.Bytes()) {
				log.Error().Msg("received bad proof from client")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				return
			}
			session := &models.Session{
				UserID:            null.IntFrom(userID),
				SessionKey:        hex.EncodeToString(k),
				SessionExpiration: time.Now().Add(time.Hour),
			}
			// save session key to database, sessions have a hard expiry of 8 hours by default, but can (and should) be ended earlier
			err = session.Insert(r.Context(), database.DB, boil.Blacklist("session_id"))
			if err != nil {
				log.Err(err).Msg("")
				ws.WriteJSON(&authenticationSocketResponse{
					Serial: -1,
				})
				return
			}
		}
	}
}

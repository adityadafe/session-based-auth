package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"session-based-auth/internal/db"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, "Server is healthy")
}

func (s *APIServer) signUpHandler(w http.ResponseWriter, r *http.Request) error {
	newUser := new(db.User)

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		return err
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	if err != nil {
		return err
	}

	newUser.Password = string(hashedPass)

	if err = s.store.CreateUser(newUser); err != nil {
		return err
	}

	fmt.Println("User created with username ->", newUser.Username)

	res := fmt.Sprintf("User Created -> %s", newUser.Username)
	WriteJson(w, http.StatusCreated, res)
	return nil
}

func (s *APIServer) signInHandler(w http.ResponseWriter, r *http.Request) error {
	newUser := new(db.User)

	if err := json.NewDecoder(r.Body).Decode(newUser); err != nil {
		return err
	}

	queriedUser, err := s.store.GetAccountByEmail(newUser.Username)
	if err != nil {
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(queriedUser.Password), []byte(newUser.Password)); err != nil {
		return err
	}

	sessionID := uuid.NewString()
	ttx := time.Now().Add(120 * time.Second)
	if err := s.store.CreateSession(sessionID, ttx); err != nil {
		return err
	}

	fmt.Println(sessionID)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  ttx,
	})

	WriteJson(w, http.StatusAccepted, "Authenticated")

	return nil
}

func (s *APIServer) signOutHandler(w http.ResponseWriter, r *http.Request) error {

	cookie, err := r.Cookie("session_id")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return fmt.Errorf("Unauthorized")
		}
		WriteJson(w, http.StatusBadRequest, "Bad request")
		return fmt.Errorf("Bad Request")
	}

	session := cookie.Value

	s.store.DeleteSession(session)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
	})

	WriteJson(w, http.StatusOK, "Successfully logout")

	return nil

}

func (s *APIServer) protectedRouteHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return err
	}

	session, err := s.store.GetSessionById(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if session.IsExpired() {
		if err := s.store.DeleteSession(session.SessionId); err != nil {
			return err
		}

		return err
	}

	s.store.DeleteSession(session.SessionId)

	sessionID := uuid.NewString()
	ttx := time.Now().Add(120 * time.Second)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  ttx,
	})

	WriteJson(w, http.StatusOK, "welcome to protected route")

	return nil
}

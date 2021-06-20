package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/haasin-farooq/ivents/api/models"
	"github.com/haasin-farooq/ivents/api/responses"
	"github.com/haasin-farooq/ivents/utils"
)

// UserSignUp controller for creating new users
func (a *App) UserSignUp(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status": "success",
		"message": "Register successfully",
	}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, _ := user.GetUser(a.DB)
	if usr != nil {
		resp["status"] = "failed"
		resp["message"] = "User already registered, please login"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	user.Prepare()

	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	userCreated, err := user.SaveUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["user"] = userCreated
	responses.JSON(w, http.StatusCreated, resp)
}

// Login signs in users
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status": "success",
		"message": "Logged in",
	}

	user := &models.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	usr, err := user.GetUser(a.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	if usr == nil {
		resp["status"] = "failed"
		resp["message"] = "Login failed, please sign up"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	err = models.CheckPasswordHash(user.Password, usr.Password)
	if err != nil {
		resp["status"] = "failed"
		resp["message"] = "Wrong password, please try again"
		responses.JSON(w, http.StatusBadRequest, resp)
		return
	}

	token, err := utils.EncodeAuthToken(usr.ID)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	resp["token"] = token
	responses.JSON(w, http.StatusOK, resp)
}
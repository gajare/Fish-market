package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gajare/Fish-market/middleware"
	"github.com/gajare/Fish-market/models"
	"github.com/gajare/Fish-market/service"
	"github.com/gajare/Fish-market/utils"
	"github.com/gorilla/mux"
)

type UserController struct {
	S *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{S: s}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var dto models.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	_, role, ok := middleware.GetUser(r)
	canSetRole := ok && role == "Admin"
	u, err := c.S.Create(r.Context(), dto, canSetRole)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	u.PasswordHash = ""
	utils.JSON(w, http.StatusCreated, u)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var dto models.LoginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}
	u, err := c.S.Login(r.Context(), dto.Email, dto.Password)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	tok, err := utils.GenerateJWT(uint(u.ID), string(u.Role))
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, "token error")
		return
	}
	utils.JSON(w, http.StatusOK, models.LoginResponse{Token: tok})
}

func (c *UserController) List(w http.ResponseWriter, r *http.Request) {
	users, err := c.S.List(r.Context())
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, users)
}

func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)

	uid, role, _ := middleware.GetUser(r)
	if role != "admin" && uid != id {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}
	u, err := c.S.GetByID(r.Context(), id)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}
	u.PasswordHash = ""
	utils.JSON(w, http.StatusOK, u)
}

func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)

	uid, role, _ := middleware.GetUser(r)
	if role != "admin" && uid != id {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}

	var dto models.UpdateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid payload")
		return
	}

	canSetRole := role == "admin"
	u, err := c.S.Update(r.Context(), id, dto, canSetRole)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	u.PasswordHash = ""
	utils.JSON(w, http.StatusOK, u)
}

func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id64, _ := strconv.ParseUint(idStr, 10, 64)
	id := uint(id64)
	_, role, _ := middleware.GetUser(r)
	if role != "admin" {
		utils.Error(w, http.StatusForbidden, "forbidden")
		return
	}
	if err := c.S.Delete(r.Context(), id); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, map[string]string{"status": "deleted"})

}

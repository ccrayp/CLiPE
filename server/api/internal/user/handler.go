package user

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repository_ *UserRepository
	debug_      bool
}

func NewUserHandler(repo *UserRepository, debug bool) *UserHandler {
	return &UserHandler{
		repository_: repo,
		debug_:      debug,
	}
}

func (h *UserHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid limit", err.Error())
		return
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid offset", err.Error())
		return
	}

	var filter UserDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	if h.debug_ {
		fmt.Printf("search user: %v\n", filter)
	}

	data, err := h.repository_.Select(&filter, limit, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	count, err := h.repository_.Count()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "error while get count", err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"users":  data,
		"limit":  limit,
		"offset": offset,
		"count":  count,
	})
}

func (h *UserHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	var dto CreateUserDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	id, err := h.repository_.Create(&dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": id,
	})
}

func (h *UserHandler) Update(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	var dto CreateUserDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid body", err.Error())
		return
	}

	err = h.repository_.Update(uint(id), &dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

func (h *UserHandler) Delete(ctx *gin.Context) {

	if auth.Require(ctx, auth.User, auth.Installer) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "invalid id", err.Error())
		return
	}

	err = h.repository_.Delete(uint(id))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

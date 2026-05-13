package request

import (
	"clipe/internal/auth"
	"clipe/pkg/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	repository_ *RequestRepository
	debug_      bool
}

func NewRequestHandler(repo *RequestRepository, debug bool) *RequestHandler {
	return &RequestHandler{
		repository_: repo,
		debug_:      debug,
	}
}

func (h *RequestHandler) Filter(ctx *gin.Context) {

	if auth.Require(ctx, auth.User) == nil {
		return
	}

	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный лимит", "Неверный лимит")
		return
	}

	offset, err := strconv.Atoi(ctx.Query("offset"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный сдвиг", "Неверный сдвиг")
		return
	}

	var filter RequestDTO
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&filter); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	data, err := h.repository_.Select(&filter, limit, offset)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при запросе данных")
		return
	}

	count, err := h.repository_.Count()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, "Ошибка при подсчёте результатов", "Ошибка при подсчёте результатов")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, gin.H{
		"requests": data,
		"limit":    limit,
		"offset":   offset,
		"count":    count,
	})
}

func (h *RequestHandler) Create(ctx *gin.Context) {

	if auth.Require(ctx, auth.DecisionServer) == nil {
		return
	}

	var dto CreateRequestDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	id, err := h.repository_.Create(&dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при создании записи")
		return
	}

	if h.debug_ {
		fmt.Printf("new request: %d\n", *id)
	}

	utils.RespondSuccess(ctx, http.StatusCreated, nil, gin.H{
		"id": id,
	})
}

func (h *RequestHandler) Update(ctx *gin.Context) {

	if auth.Require(ctx) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
		return
	}

	var dto CreateRequestDTO

	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dto); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверные параметры запроса", "Неверные параметры запроса")
		return
	}

	err = h.repository_.Update(uint(id), &dto)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при обновлении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

func (h *RequestHandler) Delete(ctx *gin.Context) {

	if auth.Require(ctx) == nil {
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, "Неверный идентификатор", "Неверный идентификатор")
		return
	}

	err = h.repository_.Delete(uint(id))
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, nil, "Ошибка при удалении записи")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, nil, nil)
}

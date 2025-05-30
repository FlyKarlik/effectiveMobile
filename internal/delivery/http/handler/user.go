package http_handler

import (
	"net/http"

	http_dto "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/handler/dto"
	http_response "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/response"
	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/FlyKarlik/effectiveMobile/internal/errs"
	_ "github.com/FlyKarlik/effectiveMobile/pkg/generics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Поиск пользователей
// @Description Возвращает список пользователей с пагинацией и фильтрацией
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param limit query int false "Лимит записей (по умолчанию 10)" default(10) minimum(1) maximum(100)
// @Param offset query int false "Смещение (по умолчанию 0)" default(0) minimum(0)
// @Param name query string false "Фильтр по имени (частичное совпадение)"
// @Param surname query string false "Фильтр по фамилии (частичное совпадение)"
// @Param patronymic query string false "Фильтр по отчеству (частичное совпадение)"
// @Param nationality query string false "Фильтр по национальности (точное совпадение)"
// @Param sex query string false "Фильтр по полу" Enums(MALE, FEMALE)
// @Param age query int false "Фильтр по возрасту (точное совпадение)" minimum(1) maximum(120)
// @Success 200 {object} generics.ItemsOutput[domain.User] "Успешный ответ"
// @Failure 400 {object} http_response.BaseResponse[any] "Невалидные параметры запроса"
// @Failure 500 {object} http_response.BaseResponse[any] "Внутренняя ошибка сервера"
// @Router /users [get]
func (h *HTTPHandler) SearchUsers(c *gin.Context) {
	pagination := http_dto.GetPaginationFromQuery(c)
	filter := http_dto.GetUserFilterFromQuery(c)

	data := h.usecase.SearchUsers(c.Request.Context(), pagination, filter)
	if !data.Success {
		http_response.New[any](c, http.StatusInternalServerError, data.Success, nil, errs.ErrUnknown)
		return
	}

	http_response.New(c, http.StatusOK, true, data, nil)
}

// @Summary Создание пользователя
// @Description Создает нового пользователя с указанными данными
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param input body domain.CreateUserInput true "Данные для создания пользователя"
// @Success 200 {object} http_response.BaseResponse[any] "Успешное создание пользователя"
// @Failure 400 {object} http_response.BaseResponse[any] "Невалидные данные запроса"
// @Failure 500 {object} http_response.BaseResponse[any] "Внутренняя ошибка сервера"
// @Router /users [post]
func (h *HTTPHandler) CreateUser(c *gin.Context) {
	var input domain.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		http_response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrInvalidRequest)
		return
	}

	if err := h.usecase.CreateUser(c.Request.Context(), input); err != nil {
		http_response.New[any](c, http.StatusInternalServerError, false, nil, errs.ErrUnknown)
		return
	}

	http_response.New[any](c, http.StatusOK, true, nil, nil)
}

// @Summary Обновление пользователя
// @Description Обновляет данные пользователя по его идентификатору
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param id path string true "UUID пользователя"
// @Param input body domain.UpdateUserInput true "Данные для обновления"
// @Success 200 {object} http_response.BaseResponse[any] "Успешное обновление"
// @Failure 400 {object} http_response.BaseResponse[any] "Невалидные параметры запроса"
// @Failure 404 {object} http_response.BaseResponse[any] "Пользователь не найден"
// @Failure 500 {object} http_response.BaseResponse[any] "Внутренняя ошибка сервера"
// @Router /users/{id} [patch]
func (h *HTTPHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		http_response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrParamRequired)
		return
	}

	var input domain.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		http_response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrInvalidRequest)
		return
	}

	if err := h.usecase.UpdateUserByID(c.Request.Context(), uuid.MustParse(id), input); err != nil {
		http_response.New[any](c, http.StatusInternalServerError, false, nil, errs.ErrUnknown)
		return
	}

	http_response.New[any](c, http.StatusOK, true, nil, nil)
}

// @Summary Удаление пользователя
// @Description Удаляет пользователя по его идентификатору
// @Tags Пользователи
// @Accept json
// @Produce json
// @Param id path string true "UUID пользователя"
// @Success 200 {object} http_response.BaseResponse[any] "Успешное удаление"
// @Failure 400 {object} http_response.BaseResponse[any] "Невалидные параметры запроса"
// @Failure 404 {object} http_response.BaseResponse[any] "Пользователь не найден"
// @Failure 500 {object} http_response.BaseResponse[any] "Внутренняя ошибка сервера"
// @Router /users/{id} [delete]
func (h *HTTPHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		http_response.New[any](c, http.StatusBadRequest, false, nil, errs.ErrParamRequired)
		return
	}

	if err := h.usecase.DeleteUserByID(c.Request.Context(), uuid.MustParse(id)); err != nil {
		http_response.New[any](c, http.StatusInternalServerError, false, nil, errs.ErrUnknown)
		return
	}

	http_response.New[any](c, http.StatusOK, true, nil, nil)
}

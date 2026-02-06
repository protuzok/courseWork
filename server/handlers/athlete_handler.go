package handlers

import (
	"courseWork/server/storage"
	"courseWork/shared"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AthleteHandler struct {
	storage *storage.AthleteRepository
}

func NewAthleteHandler(s *storage.AthleteRepository) *AthleteHandler {
	return &AthleteHandler{storage: s}
}

func (h *AthleteHandler) Create(c echo.Context) error {
	a := new(shared.Athlete)
	err := c.Bind(a)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err = h.storage.Create(c.Request().Context(), *a)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *AthleteHandler) FetchAll(c echo.Context) error {
	athletes, err := h.storage.GetAll(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, athletes)
}

func (h *AthleteHandler) Delete(c echo.Context) error {
	delReq := new(shared.DeleteRequest)

	err := c.Bind(delReq)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.storage.Delete(c.Request().Context(), delReq.IDs)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *AthleteHandler) Update(c echo.Context) error {
	a := &shared.Athlete{}
	err := c.Bind(a)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.storage.Update(c.Request().Context(), *a)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *AthleteHandler) FetchSorted(c echo.Context) error {
	athletes, err := h.storage.GetAllSortedByRun100m(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, athletes)
}

func (h *AthleteHandler) FetchBest(c echo.Context) error {
	athletes, err := h.storage.GetBestOverallAthlete(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, athletes)
}

func (h *AthleteHandler) FetchBestPressMinJump(c echo.Context) error {
	athletes, err := h.storage.GetBestPressMinJump(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, athletes)
}

func (h *AthleteHandler) FetchWithRun3kmDeviation(c echo.Context) error {
	athletes, err := h.storage.GetWithRun3kmDeviation(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, athletes)
}

func (h *AthleteHandler) FetchMinPressRun100mStats(c echo.Context) error {
	stats, err := h.storage.GetMinPressRun100mStats(c.Request().Context())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, stats)
}

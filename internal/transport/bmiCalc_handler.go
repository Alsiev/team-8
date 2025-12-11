package transport

import (
	"log/slog"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BmiHandler struct {
	log *slog.Logger
}

func NewBmiHandler(log *slog.Logger) *BmiHandler {
	return &BmiHandler{log: log}
}

func (h *BmiHandler) RegisterRoutes(r *gin.Engine) {
	bmi := r.Group("/bmi")
	{
		bmi.POST("/", h.BmiInput)
	}
}

// BmiValues godoc
// @Description Входные данные для расчёта BMI
type BmiValues struct {
	Weigth float64 `json:"weigth"`
	Heigth float64 `json:"heigth"`
}

// BmiInput godoc
// @Summary Расчёт индекса массы тела (BMI)
// @Description Принимает вес (кг) и рост (см), возвращает категорию и BMI
// @Tags BMI
// @Accept json
// @Produce json
// @Param data body BmiValues true "Данные для BMI"
// @Success 200 {object} map[string]interface{} "Успешный расчёт"
// @Failure 400 {object} map[string]string "Некорректный ввод"
// @Router /bmi/ [post]
func (h *BmiHandler) BmiInput(r *gin.Context) {
	var input BmiValues

	if err := r.ShouldBindJSON(&input); err != nil {
		h.log.Warn("erro is not valid type")
		r.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	category, result := BmiCalc(input.Weigth, input.Heigth)

	h.log.Info("operation succes")

	r.IndentedJSON(http.StatusOK, gin.H{
		"category": category,
		"result":   result,
	})
}

func BmiCalc(weight, height float64) (string, float64) {
	hMeters := height / 100
	bmi := weight / (hMeters * hMeters)
	rounded := math.Round(bmi*100) / 100
	switch {
	case rounded < 18.5:
		return "Недостаточный вес", rounded
	case rounded >= 18.5 && rounded <= 25:
		return "Нормальный вес", rounded
	default:
		return "Избыточный вес или ожирение", rounded
	}
}

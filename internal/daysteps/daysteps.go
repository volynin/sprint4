package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("некорректные данные: ожидается два значения")
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("некорректные данные: количество шагов должно быть целым числом")
	}

	// Проверка: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Преобразуем второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("некорректные данные: продолжительность должна быть в формате времени (например, 30m)")
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка: %v", err) // ВАЖНО: используем log вместо fmt.Println
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKilometers := distanceMeters / mInKm

	// Вычислим количество калорий, потраченных на прогулке
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка при вычислении калорий: %v", err) // ВАЖНО: используем log
		return ""
	}

	// Формируем строку для возврата
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKilometers,
		calories,
	)

	return result
}

package daysteps

import (
	// Добавляем необходимые пакеты
	"fmt"
	"strconv"
	"strings"
	"time"
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
		return 0, 0, err
	}

	// Проверка: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0")
	}

	// Преобразуем второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// Получаем данные о количестве шагов и продолжительности прогулки
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Проверяем, что количество шагов больше 0
	if steps <= 0 {
		return ""
	}

	// Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKilometers := distanceMeters / mInKm

	// Вычислим количество калорий, потраченных на прогулке,  WalkingSpentCalories() будет определена в пакете spentcalories
	calories := WalkingSpentCalories(weight, height, distanceKilometers, duration)

	// Формируем строку для возврата
	result := fmt.Sprintf("Пройденное расстояние: %.2f км, потраченные калории: %.2f", distanceKilometers, calories)

	return result
}

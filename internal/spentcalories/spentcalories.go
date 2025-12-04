package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("некорректные данные: ожидается три значения")
	}

	// Преобразовать первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	// Получить вид активности
	activity := parts[1]

	// Преобразуем третий элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// Рассчитываем длину шага
	stepLength := height * stepLengthCoefficient

	// Умножаем количество шагов на длину шага
	totalDistanceMeters := float64(steps) * stepLength

	// Переводим дистанцию в километры
	distanceKilometers := totalDistanceMeters / mInKm

	return distanceKilometers
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Проверяем, что продолжительность больше 0
	if duration <= 0 {
		return 0
	}

	// Вычисляем дистанцию в километрах
	dist := distance(steps, height)

	// Переводим продолжительность в часы
	durationHours := duration.Seconds() / 3600

	// Вычисляем среднюю скорость
	meanSpeed := dist / durationHours

	return meanSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// Получаем значения из строки данных
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	// Рассчитываем дистанцию, среднюю скорость и калории в зависимости от вида активности
	var distance float64
	var meanSpeed float64
	var calories float64

	switch activity {
	case "Ходьба":
		distance = distance(steps, height)
		meanSpeed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		distance = distance(steps, height)
		meanSpeed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	// Формируем строку с информацией о тренировке
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), distance, meanSpeed, calories)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные параметры на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры: все значения должны быть больше нуля")
	}

	// Рассчитываем среднюю скорость
	meanSpeed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитываем количество калорий
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// Проверяем входные параметры на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные параметры: все значения должны быть больше нуля")
	}

	// Рассчитываем среднюю скорость
	meanSpeed := meanSpeed(steps, height, duration)

	// Переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()

	// Рассчитываем количество калорий
	calories := (weight * meanSpeed * durationInMinutes) / minInH

	// Применяем корректирующий коэффициент
	calories *= walkingCaloriesCoefficient

	return calories, nil
}

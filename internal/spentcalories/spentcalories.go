package spentcalories

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// Константы для расчетов
const (
	stepLengthCoefficient       = 0.45 // коэффициент для расчета длины шага на основе роста
	mInKm                       = 1000.0
	minInH                      = 60.0
	walkingCaloriesCoefficient  = 0.035 // корректирующий коэффициент для ходьбы
	runningCaloriesCoefficient  = 0.05  // (опционально, если используется)
)

// parseTraining парсит строку формата "3456,Ходьба,3h00m"
func parseTraining(data string) (int, string, time.Duration, error) {
	// 1. Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// 2. Проверить длину слайса
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных тренировки")
	}

	// 3. Преобразовать количество шагов в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}

	activity := parts[1]

	// 4. Преобразовать продолжительность в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}

	// 5. Возврат результатов
	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	// Длина шага = рост * коэффициент
	stepLength := height * stepLengthCoefficient
	// Дистанция в метрах = шаги * длина шага
	totalMeters := float64(steps) * stepLength
	// Перевод в километры
	return totalMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. Проверить, что duration > 0
	if duration <= 0 {
		return 0
	}

	// 2. Вычислить дистанцию в км
	dist := distance(steps, height)

	// 3. Вычислить среднюю скорость (дистанция / время в часах)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}

	return dist / hours
}

// RunningSpentCalories рассчитывает калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверить входные параметры на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	// 2. Рассчитать среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитать калории: (weight * meanSpeed * durationInMinutes) / minInH
	durationInMinutes := duration.Minutes()
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

// WalkingSpentCalories рассчитывает калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверить входные параметры на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	// 2. Рассчитать среднюю скорость
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитать базовые калории в минутах
	durationInMinutes := duration.Minutes()
	baseCalories := (weight * speed * durationInMinutes) / minInH

	// 4. Умножить на коэффициент ходьбы
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

// TrainingInfo возвращает полную информацию о тренировке по строке данных
func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Получить данные из parseTraining(), логировать ошибки через log.Println(err)
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var dist float64
	var speed float64

	// 2. Определить тип тренировки через switch
	switch activity {
	case "Ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		// 4. Если тип неизвестен
		return "", errors.New("неизвестный тип тренировки")
	}

	// 3. Сформировать и вернуть строку отчета
	result := "Тип тренировки: " + activity + "\n" +
		"Длительность: " + strconv.FormatFloat(duration.Hours(), 'f', 2, 64) + " ч.\n" +
		"Дистанция: " + strconv.FormatFloat(dist, 'f', 2, 64) + " км.\n" +
		"Скорость: " + strconv.FormatFloat(speed, 'f', 2, 64) + " км/ч\n" +
		"Сожгли калорий: " + strconv.FormatFloat(calories, 'f', 2, 64)

	return result, nil
}
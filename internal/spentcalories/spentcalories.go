package spentcalories

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// Константы для расчетов по формулам Практикума
const (
	stepLengthCoefficient              = 0.45
	mInKm                              = 1000.0
	minInH                             = 60.0
	walkingCaloriesWeightMultiplier    = 0.035
	walkingSpeedHeightMultiplier       = 0.029
	runningCaloriesMeanSpeedMultiplier = 18.0
	runningCaloriesMeanSpeedShift      = 20.0
)

// parseTraining парсит строку формата "3456,Ходьба,3h00m"
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных тренировки")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	// Исправление 1: Проверка на нулевые и отрицательные шаги
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	// Исправление 2: Проверка на нулевое время
	if duration <= 0 {
		return 0, "", 0, errors.New("длительность должна быть больше нуля")
	}

	return steps, activity, duration, nil
}

// distance вычисляет дистанцию в километрах
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	totalMeters := float64(steps) * stepLength
	return totalMeters / mInKm
}

// meanSpeed вычисляет среднюю скорость в км/ч
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hours == 0 {
		return 0
	}
	return dist / hours
}

// RunningSpentCalories рассчитывает калории при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}
	speed := meanSpeed(steps, height, duration)

	// Исправление 3: Правильная формула для бега
	calories := (runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight / mInKm * duration.Minutes()
	return calories, nil
}

// WalkingSpentCalories рассчитывает калории при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}
	speed := meanSpeed(steps, height, duration)

	// Исправление 4: Правильная формула для ходьбы
	// Сначала переводим скорость из км/ч в м/с
	speedMS := speed * mInKm / (minInH * minInH)
	calories := (walkingCaloriesWeightMultiplier*weight + (speedMS*speedMS/height)*walkingSpeedHeightMultiplier*weight) * duration.Minutes()

	return calories, nil
}

// TrainingInfo возвращает полную информацию о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	var dist float64
	var speed float64

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
		return "", errors.New("неизвестный тип тренировки")
	}

	result := "Тип тренировки: " + activity + "\n" +
		"Длительность: " + strconv.FormatFloat(duration.Hours(), 'f', 2, 64) + " ч.\n" +
		"Дистанция: " + strconv.FormatFloat(dist, 'f', 2, 64) + " км.\n" +
		"Скорость: " + strconv.FormatFloat(speed, 'f', 2, 64) + " км/ч\n" +
		"Сожгли калорий: " + strconv.FormatFloat(calories, 'f', 2, 64)

	return result, nil
}

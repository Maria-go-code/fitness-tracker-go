package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	stepLengthCoefficient              = 0.65
	mInKm                              = 1000.0
	minInH                             = 60.0
	cmInM                              = 100.0
	walkingCaloriesWeightMultiplier    = 0.035
	walkingSpeedHeightMultiplier       = 0.029
	runningCaloriesMeanSpeedMultiplier = 18.0
	runningCaloriesMeanSpeedShift      = 20.0
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных тренировки")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть больше нуля")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("длительность должна быть больше нуля")
	}

	return steps, activity, duration, nil
}

func distance(steps int) float64 {
	return float64(steps) * stepLengthCoefficient / mInKm
}

func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps)
	return dist / duration.Hours()
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}
	speed := meanSpeed(steps, duration)

	calories := (runningCaloriesMeanSpeedMultiplier*speed - runningCaloriesMeanSpeedShift) * weight / mInKm * duration.Minutes()
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}

	speed := meanSpeed(steps, duration)
	speedMS := speed * mInKm / (minInH * minInH)
	heightM := height / cmInM

	calories := (walkingCaloriesWeightMultiplier*weight + (speedMS*speedMS/heightM)*walkingSpeedHeightMultiplier*weight) * duration.Minutes()
	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	var dist float64
	var speed float64

	switch activity {
	case "Ходьба":
		dist = distance(steps)
		speed = meanSpeed(steps, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		dist = distance(steps)
		speed = meanSpeed(steps, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), dist, speed, calories)

	return result, nil
}

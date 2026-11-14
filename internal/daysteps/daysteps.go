package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Константы для расчетов
const (
	stepLength = 0.65 // пример длины шага в метрах (или задается по условию проекта)
	mInKm      = 1000.0
)

// parsePackage парсит строку вида "678,0h50m" на шаги и продолжительность.
func parsePackage(data string) (int, time.Duration, error) {
	// 1. Разделить строку на слайс строк
	parts := strings.Split(data, ",")

	// 2. Проверить, чтобы длина слайса была равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// 3. Преобразовать первый элемент слайса (количество шагов) в тип int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	// 4. Проверить: количество шагов должно быть больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// 5. Преобразовать второй элемент слайса в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}

	// Дополнительная проверка на валидность duration
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть больше 0")
	}

	// 6. Возвращаем результаты
	return steps, duration, nil
}

// DayActionInfo обрабатывает строку данных и возвращает отформатированную информацию о прогулке.
func DayActionInfo(data string, weight, height float64) string {
	// 1. Получить данные через parsePackage()
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 2. Проверить, что шагов больше 0 (уже проверено в parsePackage, но продублируем для надежности)
	if steps <= 0 {
		return ""
	}

	// 3. Вычислить дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4. Перевести дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// 5. Вычислить количество калорий (используем функцию из пакета spentcalories)
	// Примечание: предполагается, что пакет spentcalories импортирован или находится в том же контексте
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 6. Сформировать и вернуть строку
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps, distanceKm, calories,
	)
}
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

// loadBankData читает banks.txt построчно и превращает строки в []Bank
func loadBankData(path string) ([]Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть %q: %w", path, err)
	}
	defer f.Close()

	var banks []Bank
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return nil, fmt.Errorf("строка %d: неверный формат (%q)", lineNum, line)
		}

		binFrom, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("строка %d: неверное BinFrom %q: %w", lineNum, parts[1], err)
		}
		binTo, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("строка %d: неверное BinTo %q: %w", lineNum, parts[2], err)
		}
		if binFrom > binTo {
			return nil, fmt.Errorf("строка %d: BinFrom больше BinTo (%d > %d)", lineNum, binFrom, binTo)
		}

		banks = append(banks, Bank{
			Name:    parts[0],
			BinFrom: binFrom,
			BinTo:   binTo,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения %q: %w", path, err)
	}

	return banks, nil
}

// getUserInput читает ввод пользователя и возвращает строку БЕЗ пробелов/дефисов.
// Пустая строка (просто Enter) — сигнал на выход.
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nВведите номер карты (Enter — выход): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка ввода:", err)
		return ""
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}

	// Уберём пробелы/дефисы, чтобы дальше всё работало с «чистыми» цифрами
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "-", "")

	return input
}

// validateInput проверяет формат: длина 13-19 и только цифры
func validateInput(cardNumber string) bool {
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}
	for i := 0; i < len(cardNumber); i++ {
		ch := cardNumber[i]
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

// validateLuhn проверяет номер по алгоритму Луна
func validateLuhn(cardNumber string) bool {
	if len(cardNumber) < 2 {
		return false
	}

	sum := 0
	double := false

	for i := len(cardNumber) - 1; i >= 0; i-- {
		ch := cardNumber[i]
		if ch < '0' || ch > '9' {
			return false
		}

		d := int(ch - '0')

		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
		double = !double
	}

	return sum%10 == 0
}

// extractBIN возвращает первые 6 цифр номера карты как int
func extractBIN(cardNumber string) int {
	if len(cardNumber) < 6 {
		return -1
	}
	bin, err := strconv.Atoi(cardNumber[:6])
	if err != nil {
		return -1
	}
	return bin
}

// identifyBank ищет банк по BIN в диапазонах
func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Неизвестный банк"
}

func main() {
	// 1) Приветствие + загрузка банков
	fmt.Println("🚀 Добро пожаловать в программу валидации карт!")

	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("❌ Ошибка загрузки данных банков:", err)
		return
	}

	fmt.Printf("✅ Загружено банков: %d\n", len(banks))

	// 2) Основной цикл
	for {
		// 3) Получение ввода
		cardNumber := getUserInput()

		// 4) Проверка на выход
		if cardNumber == "" {
			fmt.Println("👋 Программа завершена")
			break
		}

		// 5) Валидация формата
		if !validateInput(cardNumber) {
			fmt.Println("❌ Ошибка формата: номер должен содержать 13–19 цифр (без букв и символов).")
			continue
		}

		// 6) Проверка Луна
		if !validateLuhn(cardNumber) {
			fmt.Println("❌ Номер карты невалиден (не прошёл проверку Луна).")
			continue
		}

		// 7) Определение банка
		bin := extractBIN(cardNumber)
		bankName := identifyBank(bin, banks)

		// 8) Вывод результатов
		fmt.Println("✅ Номер карты валиден!")
		if bankName != "Неизвестный банк" {
			fmt.Println("🏦 Банк:", bankName)
		} else {
			fmt.Println("🏦 Эмитент не определен")
		}
	}
}

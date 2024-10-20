package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func generatePasswords() []string {
	var passwords []string
	for _, a := range letters {
		for _, b := range letters {
			for _, c := range letters {
				for _, d := range letters {
					for _, e := range letters {
						passwords = append(passwords, string([]rune{a, b, c, d, e}))
					}
				}
			}
		}
	}
	return passwords
}

func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func hashSHA256(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func bruteForceForHash(passwords []string, hash string, numThreads int, wg *sync.WaitGroup) {
	defer wg.Done()
	var localWg sync.WaitGroup
	passwordsPerThread := len(passwords) / numThreads
	found := false
	var mutex sync.Mutex

	for i := 0; i < numThreads; i++ {
		localWg.Add(1)
		go func(startIndex int) {
			defer localWg.Done()
			endIndex := startIndex + passwordsPerThread
			if endIndex > len(passwords) {
				endIndex = len(passwords)
			}
			for _, password := range passwords[startIndex:endIndex] {
				md5Hash := hashMD5(password)
				sha256Hash := hashSHA256(password)

				mutex.Lock()
				if found {
					mutex.Unlock()
					return
				}
				if md5Hash == hash || sha256Hash == hash {
					found = true
					fmt.Printf("Найден пароль для хэша %s: %s\n", hash, password)
					mutex.Unlock()
					return
				}
				mutex.Unlock()
			}
		}(i * passwordsPerThread)
	}

	localWg.Wait()
	if !found {
		fmt.Printf("Пароль для хэша %s не найден\n", hash)
	}
}

func bruteForceMultiThread(passwords []string, hashes []string, numThreads int) {
	start := time.Now()
	var wg sync.WaitGroup

	for _, hash := range hashes {
		wg.Add(1)
		go bruteForceForHash(passwords, hash, numThreads, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Время выполнения (многопоточно): %s\n", elapsed)
}

func readHashesFromFile(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	hashes := strings.Split(strings.TrimSpace(string(data)), " ")
	return hashes, nil
}

func main() {
	passwords := generatePasswords()

	predefinedHashes := []string{
		"1115dd800feaacefdf481f1f9070374a2a81e27880f187396db67958b207cbad",
		"3a7bd3e2360a3d29eea436fcfb7e44c735d117c42d1c1835420b6b9942dd4f1b",
		"74e1bb62f8dabb8125a58852b63bdf6eaef667cb56ac7f7cdba6d7305c50a22f",
		"7a68f09bd992671bb3b19a5e70b7827e",
	}

	var hashChoice int
	var hashes []string
	fmt.Println("Выберите способ ввода хэшей:")
	fmt.Println("1. Использовать заранее заданные хэши")
	fmt.Println("2. Прочитать хэши из файла")
	fmt.Scan(&hashChoice)

	switch hashChoice {
	case 1:
		hashes = predefinedHashes
		fmt.Println("Используются заранее заданные хэши.")
	case 2:
		var filename string
		fmt.Println("Введите имя файла:")
		fmt.Scan(&filename)
		var err error
		hashes, err = readHashesFromFile(filename)
		if err != nil {
			fmt.Printf("Ошибка чтения файла: %v\n", err)
			return
		}
		fmt.Println("Хэши прочитаны из файла.")
	default:
		fmt.Println("Неверный выбор.")
		return
	}

	var mode int
	fmt.Println("Выберите режим работы:")
	fmt.Println("1. Однопоточный")
	fmt.Println("2. Многопоточный")
	fmt.Scan(&mode)

	switch mode {
	case 1:
		fmt.Println("Запуск однопоточного режима:")
		start := time.Now()
		for _, hash := range hashes {
			found := false
			for _, password := range passwords {
				md5Hash := hashMD5(password)
				sha256Hash := hashSHA256(password)
				if md5Hash == hash || sha256Hash == hash {
					fmt.Printf("Найден пароль для хэша %s: %s\n", hash, password)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Пароль для хэша %s не найден\n", hash)
			}
		}
		elapsed := time.Since(start)
		fmt.Printf("Время выполнения (однопоточно): %s\n", elapsed)
	case 2:
		var numThreads int
		fmt.Println("Введите количество потоков (например, 4):")
		fmt.Scan(&numThreads)

		fmt.Println("Запуск многопоточного режима:")
		bruteForceMultiThread(passwords, hashes, numThreads)
	default:
		fmt.Println("Неверный выбор.")
	}
}

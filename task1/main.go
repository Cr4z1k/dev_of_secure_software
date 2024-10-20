package main

import (
	"archive/zip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func printDiskInfo() {
	fmt.Println("\nВыберите платформу для вывода информации о дисках:")
	fmt.Println("1. Windows")
	fmt.Println("2. macOS")

	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		if runtime.GOOS == "windows" {
			printWindowsDiskInfo()
		} else {
			fmt.Println("Информация о Windows доступна только на Windows.")
		}
	case 2:
		if runtime.GOOS == "darwin" {
			printMacDiskInfo()
		} else {
			fmt.Println("Информация о macOS доступна только на macOS.")
		}
	default:
		fmt.Println("Неверный выбор.")
	}
}

func printWindowsDiskInfo() {
	cmd := exec.Command("wmic", "logicaldisk", "get", "name,description,size,freespace")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды:", err)
		return
	}
	fmt.Println("Информация о дисках (Windows):")
	fmt.Println(string(output))
}

func printMacDiskInfo() {
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Ошибка при выполнении команды:", err)
		return
	}
	fmt.Println("Информация о дисках (macOS):")
	fmt.Println(string(output))
}

func createFile(filename string) (*os.File, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func writeFile(filename string, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func deleteFile(filename string) error {
	return os.Remove(filename)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func createJSONFile(filename string, person Person) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(person)
}

func readJSONFile(filename string) (Person, error) {
	var person Person
	data, err := os.ReadFile(filename)
	if err != nil {
		return person, err
	}
	err = json.Unmarshal(data, &person)
	return person, err
}

type Product struct {
	XMLName xml.Name `xml:"product"`
	Name    string   `xml:"name"`
	Price   float64  `xml:"price"`
}

func createXMLFile(filename string, product Product) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(product)
}

func readXMLFile(filename string) (Product, error) {
	var product Product
	data, err := os.ReadFile(filename)
	if err != nil {
		return product, err
	}
	err = xml.Unmarshal(data, &product)
	return product, err
}

func createZip(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("ошибка при создании zip файла: %v", err)
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		err = addFileToZip(zipWriter, file)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении файла %s в zip: %v", file, err)
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ошибка при открытии файла %s: %v", filename, err)
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return fmt.Errorf("не удалось получить информацию о файле %s: %v", filename, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("ошибка при создании заголовка файла %s: %v", filename, err)
	}

	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("ошибка при создании записи для файла %s: %v", filename, err)
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return fmt.Errorf("ошибка при копировании файла %s в zip: %v", filename, err)
	}

	return nil
}

func unzipFile(zipFile, destDir string) error {
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return fmt.Errorf("ошибка при открытии zip архива %s: %v", zipFile, err)
	}
	defer r.Close()

	for _, file := range r.File {
		filePath := filepath.Join(destDir, file.Name)
		fmt.Println("Разархивирование файла:", filePath)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(filePath, os.ModePerm)
			if err != nil {
				return fmt.Errorf("ошибка при создании директории %s: %v", filePath, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return fmt.Errorf("ошибка при создании директории для файла %s: %v", filePath, err)
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("ошибка при создании файла %s: %v", filePath, err)
		}

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("ошибка при чтении файла %s из архива: %v", filePath, err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("ошибка при записи содержимого файла %s: %v", filePath, err)
		}
	}
	return nil
}

func showMenu() {
	fmt.Println("\nВыберите действие:")
	fmt.Println("1. Показать информацию о дисках (Windows/macOS)")
	fmt.Println("2. Работа с файлами")
	fmt.Println("3. Работа с JSON")
	fmt.Println("4. Работа с XML")
	fmt.Println("5. Работа с ZIP архивом")
	fmt.Println("6. Выход")
}

func main() {
	for {
		showMenu()

		var choice int
		fmt.Print("\nВведите номер действия: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			printDiskInfo()

		case 2:
			fmt.Println("\nВыбрана работа с файлами.")
			var filename string
			fmt.Print("Введите имя файла: ")
			fmt.Scan(&filename)

			fmt.Println("1. Создать файл")
			fmt.Println("2. Записать в файл строку")
			fmt.Println("3. Прочитать файл")
			fmt.Println("4. Удалить файл")
			fmt.Print("\nВведите действие с файлом: ")
			var fileAction int
			fmt.Scan(&fileAction)

			switch fileAction {
			case 1:
				createFile(filename)
				fmt.Println("Файл создан.")
			case 2:
				fmt.Print("Введите строку для записи в файл: ")
				var content string
				fmt.Scan(&content)
				writeFile(filename, content)
				fmt.Println("Запись завершена.")
			case 3:
				content, _ := readFile(filename)
				fmt.Println("Содержимое файла:", content)
			case 4:
				deleteFile(filename)
				fmt.Println("Файл удалён.")
			}

		case 3:
			fmt.Println("\nВыбрана работа с JSON.")
			var filename string
			fmt.Print("Введите имя JSON файла: ")
			fmt.Scan(&filename)

			fmt.Println("1. Создать и записать объект в JSON")
			fmt.Println("2. Прочитать объект из JSON")
			fmt.Println("3. Удалить JSON файл")
			fmt.Print("\nВведите действие с JSON: ")
			var jsonAction int
			fmt.Scan(&jsonAction)

			switch jsonAction {
			case 1:
				var name string
				var age int
				fmt.Print("Введите имя: ")
				fmt.Scan(&name)
				fmt.Print("Введите возраст: ")
				fmt.Scan(&age)
				person := Person{Name: name, Age: age}
				createJSONFile(filename, person)
				fmt.Println("JSON файл создан и данные записаны.")
			case 2:
				person, _ := readJSONFile(filename)
				fmt.Printf("Имя: %s, Возраст: %d\n", person.Name, person.Age)
			case 3:
				deleteFile(filename)
				fmt.Println("JSON файл удалён.")
			}

		case 4:
			fmt.Println("\nВыбрана работа с XML.")
			var filename string
			fmt.Print("Введите имя XML файла: ")
			fmt.Scan(&filename)

			fmt.Println("1. Создать и записать объект в XML")
			fmt.Println("2. Прочитать объект из XML")
			fmt.Println("3. Удалить XML файл")
			fmt.Print("\nВведите действие с XML: ")
			var xmlAction int
			fmt.Scan(&xmlAction)

			switch xmlAction {
			case 1:
				var name string
				var price float64
				fmt.Print("Введите название продукта: ")
				fmt.Scan(&name)
				fmt.Print("Введите цену продукта: ")
				fmt.Scan(&price)
				product := Product{Name: name, Price: price}
				createXMLFile(filename, product)
				fmt.Println("XML файл создан и данные записаны.")
			case 2:
				product, _ := readXMLFile(filename)
				fmt.Printf("Продукт: %s, Цена: %.2f\n", product.Name, product.Price)
			case 3:
				deleteFile(filename)
				fmt.Println("XML файл удалён.")
			}

		case 5:
			fmt.Println("\nВыбрана работа с ZIP архивами.")
			fmt.Println("1. Создать ZIP архив")
			fmt.Println("2. Разархивировать ZIP архив")
			fmt.Println("3. Удалить ZIP архив")
			fmt.Print("\nВведите действие с ZIP: ")
			var zipAction int
			fmt.Scan(&zipAction)

			switch zipAction {
			case 1:
				var zipFilename string
				fmt.Print("Введите имя ZIP архива: ")
				fmt.Scan(&zipFilename)

				var files string
				fmt.Print("Введите файлы для добавления через запятую: ")
				fmt.Scan(&files)
				fileList := strings.Split(files, ",")

				err := createZip(zipFilename, fileList)
				if err != nil {
					fmt.Println("Ошибка при создании архива:", err)
				} else {
					fmt.Println("ZIP архив создан.")
				}
			case 2:
				var zipFilename string
				fmt.Print("Введите имя ZIP архива: ")
				fmt.Scan(&zipFilename)

				var destDir string
				fmt.Print("Введите папку для разархивирования: ")
				fmt.Scan(&destDir)

				err := unzipFile(zipFilename, destDir)
				if err != nil {
					fmt.Println("Ошибка при разархивировании:", err)
				} else {
					fmt.Println("ZIP архив разархивирован.")
				}
			case 3:
				var zipFilename string
				fmt.Print("Введите имя ZIP архива для удаления: ")
				fmt.Scan(&zipFilename)
				deleteFile(zipFilename)
				fmt.Println("ZIP файл удалён.")
			}

		case 6:
			fmt.Println("Выход из программы.")
			return
		}
	}
}

import os
import json
import xml.etree.ElementTree as ET
import zipfile
import platform

def print_disk_info():
    print("\nВыберите платформу для вывода информации о дисках:")
    print("1. Windows")
    print("2. macOS")

    choice = input("Введите номер выбора: ")

    if choice == "1":
        if platform.system() == "Windows":
            os.system("wmic logicaldisk get name,description,size,freespace")
        else:
            print("Информация о Windows доступна только на Windows.")
    elif choice == "2":
        if platform.system() == "Darwin":
            os.system("df -h")
        else:
            print("Информация о macOS доступна только на macOS.")
    else:
        print("Неверный выбор.")

def file_operations():
    print("\nВыбрана работа с файлами.")
    filename = input("Введите имя файла: ")

    print("1. Создать файл")
    print("2. Записать в файл строку")
    print("3. Прочитать файл")
    print("4. Удалить файл")
    file_action = input("Введите действие с файлом: ")

    if file_action == "1":
        with open(filename, 'w') as f:
            pass
        print("Файл создан.")
    elif file_action == "2":
        content = input("Введите строку для записи в файл: ")
        with open(filename, 'w') as f:
            f.write(content)
        print("Запись завершена.")
    elif file_action == "3":
        if os.path.isfile(filename):
            with open(filename, 'r') as f:
                content = f.read()
            print("Содержимое файла:", content)
        else:
            print("Файл не найден.")
    elif file_action == "4":
        if os.path.isfile(filename):
            os.remove(filename)
            print("Файл удалён.")
        else:
            print("Файл не найден.")

def json_operations():
    print("\nВыбрана работа с JSON.")
    filename = input("Введите имя JSON файла: ")

    print("1. Создать и записать объект в JSON")
    print("2. Прочитать объект из JSON")
    print("3. Удалить JSON файл")
    json_action = input("Введите действие с JSON: ")

    if json_action == "1":
        name = input("Введите имя: ")
        age = input("Введите возраст: ")
        person = {"name": name, "age": int(age)}
        with open(filename, 'w') as f:
            json.dump(person, f)
        print("JSON файл создан и данные записаны.")
    elif json_action == "2":
        if os.path.isfile(filename):
            with open(filename, 'r') as f:
                person = json.load(f)
            print(f"Имя: {person['name']}, Возраст: {person['age']}")
        else:
            print("JSON файл не найден.")
    elif json_action == "3":
        if os.path.isfile(filename):
            os.remove(filename)
            print("JSON файл удалён.")
        else:
            print("JSON файл не найден.")

def xml_operations():
    print("\nВыбрана работа с XML.")
    filename = input("Введите имя XML файла: ")

    print("1. Создать и записать объект в XML")
    print("2. Прочитать объект из XML")
    print("3. Удалить XML файл")
    xml_action = input("Введите действие с XML: ")

    if xml_action == "1":
        name = input("Введите название продукта: ")
        price = input("Введите цену продукта: ")
        product = ET.Element("product")
        name_elem = ET.SubElement(product, "name")
        name_elem.text = name
        price_elem = ET.SubElement(product, "price")
        price_elem.text = price
        tree = ET.ElementTree(product)
        tree.write(filename)
        print("XML файл создан и данные записаны.")
    elif xml_action == "2":
        if os.path.isfile(filename):
            tree = ET.parse(filename)
            root = tree.getroot()
            name = root.find("name").text
            price = root.find("price").text
            print(f"Продукт: {name}, Цена: {price}")
        else:
            print("XML файл не найден.")
    elif xml_action == "3":
        if os.path.isfile(filename):
            os.remove(filename)
            print("XML файл удалён.")
        else:
            print("XML файл не найден.")

def zip_operations():
    print("\nВыбрана работа с ZIP архивами.")
    print("1. Создать ZIP архив")
    print("2. Разархивировать ZIP архив")
    print("3. Удалить ZIP архив")
    zip_action = input("Введите действие с ZIP: ")

    if zip_action == "1":
        zip_filename = input("Введите имя ZIP архива: ")
        files = input("Введите файлы для добавления через запятую: ").split(',')
        with zipfile.ZipFile(zip_filename, 'w') as zipf:
            for file in files:
                zipf.write(file.strip())
        print("ZIP архив создан.")
    elif zip_action == "2":
        zip_filename = input("Введите имя ZIP архива: ")
        dest_dir = input("Введите папку для разархивирования: ")
        with zipfile.ZipFile(zip_filename, 'r') as zipf:
            zipf.extractall(dest_dir)
        print("ZIP архив разархивирован.")
    elif zip_action == "3":
        zip_filename = input("Введите имя ZIP архива для удаления: ")
        if os.path.isfile(zip_filename):
            os.remove(zip_filename)
            print("ZIP файл удалён.")
        else:
            print("ZIP файл не найден.")

def show_menu():
    print("\nВыберите действие:")
    print("1. Показать информацию о дисках (Windows/macOS)")
    print("2. Работа с файлами")
    print("3. Работа с JSON")
    print("4. Работа с XML")
    print("5. Работа с ZIP архивом")
    print("6. Выход")

def main():
    while True:
        show_menu()
        choice = input("Введите номер действия: ")

        if choice == "1":
            print_disk_info()
        elif choice == "2":
            file_operations()
        elif choice == "3":
            json_operations()
        elif choice == "4":
            xml_operations()
        elif choice == "5":
            zip_operations()
        elif choice == "6":
            print("Выход из программы.")
            break
        else:
            print("Неверный выбор.")

if __name__ == "__main__":
    main()

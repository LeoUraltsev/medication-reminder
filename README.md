# MedicationReminder

---

### Запуск приложения:

Перед запуском приложения нужно подготовить конфиг и .env файл.   
Пример конфига лежит по пути `.config/example_config.yml`   
Пример env файла лежит в корне проета `.env`

Запустить приложение можно двумя способами:

1. Указав путь в флаг `-config-path`
    - Вручную
       ```bash
             go build -o ./build/mreminder ./cmd/medication-reminder/main.go
             ./build/mreminder -config-path=./config/config.yml
       ```
    - Через make предварительно указав путь до конфига в `Makefile`
         ```bash
         make run_app_cfgpath
         ```

2. Указав путь в переменной окружения `MR_CONFIG_PATH` или в `.env`
    - Вручную
         ```bash
         go build -o ./build/mreminder ./cmd/medication-reminder/main.go
         ./build/mreminder
         ```
    - Через make
         ```bash
         make run_app
         ```
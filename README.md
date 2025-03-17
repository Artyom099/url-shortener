# url-shortener

## todo

### 1. Как запустить проект:
```shell
go run main.go
```

### 2. Как запустить тесты?

Сейчас если configPath задавать так: 
```go
configPath := os.Getenv("CONFIG_PATH")
```
то падает такая ошибка:
```shell
level=ERROR msg="failed to initialize storage" env=local
error="storage.sqlite.NewStorage: unable to open database file: no such file or directory"
```
надо починить
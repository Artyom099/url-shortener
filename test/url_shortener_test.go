package test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"net/url"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/random"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	// Универсальный способ создания URL
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	// Создаем клиент httpexpect
	e := httpexpect.Default(t, u.String())

	e.POST("/url"). // Отправляем POST-запрос, путь - '/url'
			WithJSON(save.Request{ // Формируем тело запроса
			URL:   gofakeit.URL(),             // Генерируем случайный URL
			Alias: random.NewRandomString(10), // Генерируем случайную строку
		}).
		WithBasicAuth("myuser", "mypass"). // Добавляем к запросу креды авторизации
		Expect().                          // Далее перечисляем наши ожидания от ответа
		Status(200).                       // Код должен быть 200
		JSON().Object().                   // Получаем JSON-объект тела ответа
		ContainsKey("alias")               // Проверяем, что в нём есть ключ 'alias'
}

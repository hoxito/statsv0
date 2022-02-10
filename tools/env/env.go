package env

import (
	"os"
	"strconv"
)

// Configuration properties
type Configuration struct {
	Port              int    `json:"port"`
	RabbitURL         string `json:"rabbitUrl"`
	MongoURL          string `json:"mongoUrl"`
	RedisURL          string `json:"redisUrl"`
	WWWWPath          string `json:"wwwPath"`
	JWTSecret         string `json:"jwtSecret"`
	SecurityServerURL string `json:"securityServerUrl"`
	CatalogURL        string `json:"CatalogURL"`
	OrdersURL         string `json:"OrdersURL"`
}

var config *Configuration

func new() *Configuration {
	return &Configuration{
		Port:              3010,
		RabbitURL:         "amqp://localhost",
		MongoURL:          "mongodb+srv://jose:statsgo@cluster0.j1j5b.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
		RedisURL:          "localhost:6379",
		WWWWPath:          "www",
		JWTSecret:         "ecb6d3479ac3823f1da7f314d871989b",
		SecurityServerURL: "http://localhost:3000",
		CatalogURL:        "http://localhost:3002",
		OrdersURL:         "http://localhost:3004",
	}
}

// Get Obtiene las variables de entorno del sistema
func Get() *Configuration {
	if config == nil {
		config = load()
	}

	return config
}

// Load file properties
func load() *Configuration {
	result := new()

	if value := os.Getenv("REDIS_URL"); len(value) > 0 {
		result.RedisURL = value
	}
	if value := os.Getenv("RABBIT_URL"); len(value) > 0 {
		result.RabbitURL = value
	}

	if value := os.Getenv("MONGO_URL"); len(value) > 0 {
		result.MongoURL = value
	}

	if value := os.Getenv("PORT"); len(value) > 0 {
		if intVal, err := strconv.Atoi(value); err != nil {
			result.Port = intVal
		}
	}
	if value := os.Getenv("AUTH"); len(value) > 0 {
		result.SecurityServerURL = value
	}

	if value := os.Getenv("CATALOG"); len(value) > 0 {
		result.CatalogURL = value
	}

	if value := os.Getenv("ORDERS"); len(value) > 0 {
		result.OrdersURL = value
	}
	if value := os.Getenv("WWW_PATH"); len(value) > 0 {
		result.WWWWPath = value
	}

	if value := os.Getenv("JWT_SECRET"); len(value) > 0 {
		result.JWTSecret = value
	}

	return result
}

// @title SIP Bot API
// @version 1.0
// @description API для получения постов с JSONPlaceholder
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import "go-service/internal/application"

func main() {
	application.NewApplication().Run()
}

init_dependency:
	@go get -u github.com/antonfisher/nested-logrus-formatter
	@go get -u github.com/gin-gonic/gin
	@go get -u golang.org/x/crypto
	@go get -u gorm.io/gorm
	@go get -u gorm.io/driver/postgres
	@go get -u github.com/sirupsen/logrus
	@go get -u github.com/joho/godotenv

copy_env:
	cp .env .env.local

run:
	@go run cmd/main.go
module github.com/my-username/billion-user-app/services/media-service

go 1.21.0

require (
	github.com/gofiber/fiber/v2 v2.52.0
	github.com/my-username/billion-user-app/pkg/config v0.0.0
	github.com/my-username/billion-user-app/pkg/database v0.0.0
	github.com/my-username/billion-user-app/pkg/jwtutils v0.0.0
	github.com/my-username/billion-user-app/pkg/logger v0.0.0
	gorm.io/gorm v1.25.5
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gorm.io/driver/postgres v1.5.4 // indirect
)

replace (
	github.com/my-username/billion-user-app/pkg/config => ../../pkg/config
	github.com/my-username/billion-user-app/pkg/database => ../../pkg/database
	github.com/my-username/billion-user-app/pkg/jwtutils => ../../pkg/jwtutils
	github.com/my-username/billion-user-app/pkg/logger => ../../pkg/logger
)

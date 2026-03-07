# аналитик решил не чинить подгрузку из .env, а мне впадлу теперь


export API_PORT=8080
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=devspace
export DB_SSLMODE=disable
export GORM_LOGGER=false
export ARGON_MEMORY_USAGE=54
export ARGON_ITERATIONS=5
export ARGON_PARALL=4
export ARGON_SALT_LEN=16
export ARGON_KEY_LEN=32


cd cmd/api/ && go run .

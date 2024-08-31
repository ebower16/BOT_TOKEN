package systems

import (
	"github.com/joho/godotenv"
	"os"
	//"os"
)

func BotToken() string {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}
	return os.Getenv("BOT_TOKEN")
}

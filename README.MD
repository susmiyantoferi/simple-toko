Teknologi Yang Di Gunakan :

Go Gin Gonic
GORM
Mysql
ci/cd
godotenv

Set up local :
clone repo
ubah nama file .env.sample jadi .env
ubah file .env sesuai dengan nama database, username, password, host dan port masing masing
buat database sesuai dengan nama database yang ada di .env
un-comment code berikut di file config/db.go :
// "github.com/joho/godotenv"
// err := godotenv.Load()
	// if err != nil{
	// 	log.Fatal("error load env")
	// }
jalankan di terminal go run main.go
URL Documentation API : https://documenter.getpostman.com/view/22397647/2sB3BLjnnp
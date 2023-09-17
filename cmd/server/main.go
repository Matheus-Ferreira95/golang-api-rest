package main

import (
	"awesomeProject/configs"
	_ "awesomeProject/docs"
	"awesomeProject/internal/entity"
	"awesomeProject/internal/infra/database"
	"awesomeProject/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// @title            Go Expert API Example
// @version          1.0
// @description	     Product API with authentication
// @termsOfService   http://swagger.io/terms/
//
// @contact.name	  Wesley Willians
// @contact.url	      http://www.fullcycle.com.br
// @contact.email     atendimento@fullcycle.com.br
//
// @license.name	  Full Cycle License
// @license.url       http://www.fullcycle.com.br
//
// @host	          localhost:8000
// @BasePath          /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	var configs, err = configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dsn := "root:root@tcp(localhost:3306)/goexpert"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", configs.JWTExperisIn))

	r.Route("/products", func(route chi.Router) {
		route.Use(jwtauth.Verifier(configs.TokenAuth)) // sempre que acessar um dos serviços de products, ele vai buscar o token que pode estar na url, no body, só vai recuperar o token, nao vai validar nesse momento
		// se não achar token, retornara 401 com msg token not found
		route.Use(jwtauth.Authenticator) // vai validar se o token não é expirado e se a assinatura foi gerada por nossa aplicação, se o token for invalido retorna 401
		route.Post("/", productHandler.CreateProduct)
		route.Get("/", productHandler.GetProducts)
		route.Get("/{id}", productHandler.GetProduct)
		route.Put("/{id}", productHandler.UpdateProduct)
		route.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(route chi.Router) {
		route.Post("/", userHandler.Create)
		route.Post("/generate_token", userHandler.GetJWT)
	})

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))

	http.ListenAndServe(":8000", r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		log.Println("kkkkkk")
		next.ServeHTTP(w, r)
	})
}

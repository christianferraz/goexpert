package main

import (
	"net/http"

	"github.com/christianferraz/goexpert/9-APIs/configs"
	_ "github.com/christianferraz/goexpert/9-APIs/docs"
	"github.com/christianferraz/goexpert/9-APIs/internal/entity"
	"github.com/christianferraz/goexpert/9-APIs/internal/infra/database"
	"github.com/christianferraz/goexpert/9-APIs/internal/webserver/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title           Go Expert API Example
// @version         1.0
// @description     Product API with auhtentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   Christian Ferraz
// @contact.url    http://www.fullcycle.com.br
// @contact.email  atendimento@fullcycle.com.br

// @license.name   Full Cycle License
// @license.url    http://www.fullcycle.com.br

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})
	productDB := database.NewProductDB(db)
	userDB := database.NewUserDB(db)
	productHandler := handler.NewProductHandler(productDB)
	userHandler := handler.NewUserHandler(userDB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// recupera a aplicaçao se der um bug
	r.Use(middleware.Recoverer)

	r.Route("/products", func(r chi.Router) {
		// adicionar o middleware para pegar os dados depois
		// jwauth.Verifier vai verificar se existe algum token no request, pode ser na url, no header, etc
		// no jwauth.Authenticator vai passar a assinatura do token para posteriormente validar
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		// aqui ele valida o token com as configurações passadas acima se o token estiver invalido, vai
		// retornar token is unauthorized com erro 401
		// se estiver válido vai retornar o token status 200 e vai liberar os comandos abaixo restantes
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})
	r.Route("/users", func(r chi.Router) {
		// injeta o valor do de jwt no request e dará pra recuperar com r.context().Value("jwt")
		r.Use(middleware.WithValue("jwt", configs.TokenAuth))
		r.Use(middleware.WithValue("jwtExpiresIn", configs.JWTExpiresIn))
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWTInput)

	})
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/doc.json")))
	http.ListenAndServe(":8000", r)
}

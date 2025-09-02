package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/christianferraz/goexpert/9-APIs/internal/dto"
	"github.com/christianferraz/goexpert/9-APIs/internal/entity"
	"github.com/christianferraz/goexpert/9-APIs/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type Error struct {
	Message string `json:"message"`
}

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(db database.UserInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

// GetJWT godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request   body     dto.GetJWTInput  true  "user credentials"
// @Success      200  {object}  dto.GetJWTOutput
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /users/generate_token [post]
func (h *UserHandler) GetJWTInput(w http.ResponseWriter, r *http.Request) {
	// faz um casting pra forçar que é um (*jwtauth.JWTAuth)
	// recupera os dados do middleware do chi em main.go de
	// r.Use(middleware.WithValue("jwt", configs.TokenAuth))
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)
	var user dto.GetJWTInput
	//pega os dados passados na autenticação
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := h.UserDB.FindByEmail(user.Email)
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if !u.ValidatePassword(user.Password) {
		http.Error(w, "erro", http.StatusUnauthorized)
		return
	}
	// criar o token com o mapa de chave string e valor de qualquer coisa
	_, tokenString, err := jwt.Encode(map[string]any{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//struct anonima
	// accessToken := struct {
	// 	AccessToken string `json:"access_token"`
	// }{
	// 	AccessToken: tokenString,
	// }
	accessToken := dto.GetJWTOutput{AccessToken: tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}

// Create user godoc
// @Summary      Create user
// @Description  Create user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request     body      dto.CreateUserInput  true  "user request"
// @Success      201
// @Failure      500         {object}  Error
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.UserDB.Create(u)
	if err != nil {
		error := Error{Message: err.Error()}
		json.NewEncoder(w).Encode(error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

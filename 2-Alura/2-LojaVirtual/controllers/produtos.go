package controllers

import (
	"net/http"
	"strconv"
	"text/template"

	"github.com/christianferraz/goexpert/2-Alura/2-LojaVirtual/db"
	"github.com/christianferraz/goexpert/2-Alura/2-LojaVirtual/models"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var funcMap = template.FuncMap{
	"formatarPreco": formatarPreco,
}

var temp = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

func formatarPreco(preco float64) string {
	p := message.NewPrinter(language.BrazilianPortuguese)
	return p.Sprintf("R$ %.2f", preco)
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := db.ConectaComBancoDeDados()
	defer db.Close()

	res, err := db.Query("select * from produtos")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Close()

	pd := models.NovoProduto(db)
	produtos, err := pd.ListarTodosProdutos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = temp.ExecuteTemplate(w, "Index", produtos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func New(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nome := r.FormValue("nome")
		descricao := r.FormValue("descricao")
		preco := r.FormValue("preco")
		quantidade := r.FormValue("quantidade")

		db := db.ConectaComBancoDeDados()
		defer db.Close()

		pd := models.NovoProduto(db)
		precoConvertido, err := strconv.ParseFloat(preco, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		quantidadeConvertida, err := strconv.Atoi(quantidade)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pd.CriarNovoProduto(nome, descricao, precoConvertido, quantidadeConvertida)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	idProduto := r.URL.Query().Get("id")
	db := db.ConectaComBancoDeDados()
	defer db.Close()

	pd := models.NovoProduto(db)
	err := pd.DeletarProduto(idProduto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	idProduto := r.URL.Query().Get("id")
	db := db.ConectaComBancoDeDados()
	defer db.Close()

	pd := models.NovoProduto(db)
	produto, err := pd.BuscarProduto(idProduto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	temp.ExecuteTemplate(w, "Edit", produto)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		nome := r.FormValue("nome")
		descricao := r.FormValue("descricao")
		preco := r.FormValue("preco")
		quantidade := r.FormValue("quantidade")

		db := db.ConectaComBancoDeDados()
		defer db.Close()

		pd := models.NovoProduto(db)
		idConvertido, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		precoConvertido, err := strconv.ParseFloat(preco, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		quantidadeConvertida, err := strconv.Atoi(quantidade)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = pd.AtualizarProduto(&idConvertido, &nome, &descricao, &precoConvertido, &quantidadeConvertida)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

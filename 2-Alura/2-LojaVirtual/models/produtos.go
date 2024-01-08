package models

import "database/sql"

type Produto struct {
	Id         int
	Nome       string
	Descricao  string
	Preco      float64
	Quantidade int
	DB         *sql.DB
}

func NovoProduto(db *sql.DB) *Produto {
	return &Produto{
		DB: db,
	}
}

func (p *Produto) AtualizarProduto(id *int, nome, descricao *string, preco *float64, quantidade *int) error {
	stmt, err := p.DB.Prepare("update produtos set nome = $1, descricao = $2, preco = $3, quantidade = $4 where id = $5")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nome, descricao, preco, quantidade, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *Produto) BuscarProduto(id string) (Produto, error) {
	var produto Produto
	res, err := p.DB.Query("select * from produtos where id = $1", id)
	if err != nil {
		return produto, err
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&produto.Id, &produto.Nome, &produto.Descricao, &produto.Preco, &produto.Quantidade)
		if err != nil {
			return produto, err
		}
	}
	return produto, nil
}

func (p *Produto) ListarTodosProdutos() ([]Produto, error) {
	res, err := p.DB.Query("select * from produtos")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var produtos []Produto
	for res.Next() {
		var produto Produto
		err = res.Scan(&produto.Id, &produto.Nome, &produto.Descricao, &produto.Preco, &produto.Quantidade)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, produto)
	}
	return produtos, nil
}

func (p *Produto) CriarNovoProduto(nome, descricao string, preco float64, quantidade int) error {
	stmt, err := p.DB.Prepare("insert into produtos (nome, descricao, preco, quantidade) values ($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(nome, descricao, preco, quantidade)
	if err != nil {
		return err
	}
	return nil
}

func (p *Produto) DeletarProduto(id string) error {
	stmt, err := p.DB.Prepare("delete from produtos where id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

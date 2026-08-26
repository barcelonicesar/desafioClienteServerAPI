package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type Cambio struct {
	Bid string `json:"bid"`
}

type RespostaAPI struct {
	USDBRL Cambio `json:"USDBRL"`
}

func main() {
	var err error

	db, err = sql.Open("sqlite", "cotacoes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/cotacao", BuscaCotacaoHandler)

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func BuscaCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	// ctx, cancel := context.WithTimeout(r.Context(), 200*time.Second)
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	cambio, err := BuscaCotacao(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout de 200ms ao consultar a API")
		}

		http.Error(w, "Erro ao consultar a cotação", http.StatusInternalServerError)
		return
	}

	// ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Second)
	ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancelDB()

	_, err = db.ExecContext(
		ctxDB,
		"INSERT INTO cotacoes (bid) VALUES (?)",
		cambio.Bid,
	)
	if err != nil {
		if ctxDB.Err() == context.DeadlineExceeded {
			log.Println("Timeout de 10ms ao salvar no banco")
		}

		http.Error(w, "Erro ao salvar a cotação", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cambio)
}

func BuscaCotacao(ctx context.Context) (*Cambio, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://economia.awesomeapi.com.br/json/last/USD-BRL",
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var resposta RespostaAPI
	err = json.Unmarshal(body, &resposta)
	if err != nil {
		return nil, err
	}

	return &resposta.USDBRL, nil
}

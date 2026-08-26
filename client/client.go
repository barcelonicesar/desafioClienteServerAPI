package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cambio struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		300*time.Millisecond,
	)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"http://localhost:8080/cotacao",
		nil,
	)
	if err != nil {
		log.Println("Erro ao criar a requisição:", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout de 300ms ao consultar o servidor")
			return
		}

		log.Println("Erro ao consultar o servidor:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf(
			"Servidor retornou o status %d: %s",
			resp.StatusCode,
			string(body),
		)
		return
	}

	var cambio Cambio

	err = json.NewDecoder(resp.Body).Decode(&cambio)
	if err != nil {
		log.Println("Erro ao ler a resposta do servidor:", err)
		return
	}

	conteudo := fmt.Sprintf("Dólar: %s", cambio.Bid)

	err = os.WriteFile("cotacao.txt", []byte(conteudo), 0644)
	if err != nil {
		log.Println("Erro ao salvar o arquivo:", err)
		return
	}

	fmt.Println("Cotação salva com sucesso:", conteudo)
}

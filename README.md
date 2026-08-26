# Desafio Client-Server API

Este projeto foi desenvolvido como parte do curso Go Expert.

O objetivo é praticar os seguintes conceitos da linguagem Go:

- Servidor HTTP
- Cliente HTTP
- Context e timeout
- Consumo de API externa
- Banco de dados SQLite
- Manipulação de arquivos
- JSON

## Funcionamento

O projeto possui dois programas:

- `server`: consulta a cotação do dólar na AwesomeAPI, salva o valor no SQLite e retorna a cotação em JSON.
- `client`: consulta o servidor local e salva o valor recebido no arquivo `cotacao.txt`.

## Estrutura

desafio/
├── client/
│   └── client.go
├── server/
│   └── server.go
├── go.mod
└── go.sum
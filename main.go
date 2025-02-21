package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Estrutura resposta da API do ViaCEP
type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

// Estrutura resposta da API do BrasilAPI
type BrasilApiResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func buscaCep(url string, ch chan<- string, apiName string) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ch <- fmt.Sprintf("Erro ao criar requisição para %s: %v", apiName, err)
		return
	}

	// As informações abaixo foram utilizadas para testar as três hipóteses previstas
	// if apiName == "ViaCEP" {
	//  	time.Sleep(time.Millisecond * 750)
	// }
	// time.Sleep(time.Millisecond * 1100)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("Erro ao acessar %s: %v", apiName, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("Erro ao ler resposta de %s: %v", apiName, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("%s retornou status code %d", apiName, resp.StatusCode)
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		ch <- fmt.Sprintf("Erro ao decodificar JSON de %s: %v", apiName, err)
		return
	}

	ch <- fmt.Sprintf("Resposta de %s: %+v", apiName, result)
}

func main() {

	cep := "49142442"
	chViaCep := make(chan string)
	chBrasilApi := make(chan string)

	// URLs das APIs
	viaCepURL := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	brasilApiURL := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)

	// Goroutines para buscar os dados das APIs
	go buscaCep(viaCepURL, chViaCep, "ViaCEP")
	go buscaCep(brasilApiURL, chBrasilApi, "BrasilAPI")

	select {
	case msgViaCep := <-chViaCep:
		fmt.Println(msgViaCep)
	case msgBrasilCep := <-chBrasilApi:
		fmt.Println(msgBrasilCep)
	case <-time.After(time.Second * 1):
		fmt.Println("Erro: Timeout, excedeu 1 segundo")
	}
}

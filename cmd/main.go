package main

import (
	"bufio"
	"fmt"
	"github.com/rafaelcmd/fc-goexpert-multithreading/internal/api"
	"github.com/rafaelcmd/fc-goexpert-multithreading/internal/service"
	"os"
	"strings"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Digite o CEP (somente números) e pressione Enter: ")
	zipCode, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	zipCode = strings.TrimSpace(zipCode)

	apiURLs := []string{
		"https://brasilapi.com.br/api/cep/v1/%s",
		"http://viacep.com.br/ws/%s/json/",
	}

	httpClient := api.NewClient(1 * time.Second)

	checker := service.NewCheckerService(httpClient, apiURLs)

	result, url, err := checker.CheckZipCode(zipCode)
	if err != nil {
		fmt.Println("Error checking ZIP Code:", err)
		return
	}

	var fromAPI string

	if strings.Contains(url, "brasilapi") {
		fromAPI = "BrasilAPI"
	} else {
		fromAPI = "ViaCEP"
	}

	fmt.Println("ZIP Code data:", result, "from API:", fromAPI)
}

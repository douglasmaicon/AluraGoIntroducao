package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const MONITORAMENTO = 5
const DELAY = 3

func main() {
	exibirIntroducao()
	for {
		exibirMenu()
		comando := lerComando()
		switch comando {
		case 0:
			fmt.Println("Saindo do sistema...")
			os.Exit(0)
		case 1:
			iniciarMonitoramento()
		case 2:
			fmt.Println("Exibindo LOGs...")
			exibirLogs()
		default:
			fmt.Println("Não conheçlo este comando...")
			os.Exit(-1)
		}
	}

}

func exibirIntroducao() {
	nome := "Douglas"
	fmt.Printf("Olá sr., %s\n", nome)
}

func exibirMenu() {
	fmt.Println("1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair")
}

func lerComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitorando...")
	sites := lerArquivo()
	for i := 0; i < MONITORAMENTO; i++ {
		for _, site := range sites {
			fmt.Println("Testando o site:", site)
			testarSite(site)
		}
		time.Sleep(time.Second * DELAY)
	}

}

func lerArquivo() []string {
	var conteudo []string
	arquivo, err := os.Open("sites.txt")
	defer arquivo.Close()
	if err != nil {
		log.Fatalln("Ocorreu um erro ao abrir o arquivo de sites:", err)
	} else {
		leitor := bufio.NewReader(arquivo)
		for {
			linha, err := leitor.ReadString('\n')
			linha = strings.TrimSpace(linha)
			if err == io.EOF {
				break
			} else {
				fmt.Println(linha)
				conteudo = append(conteudo, linha)
			}
		}
	}

	return conteudo
}

func testarSite(site string) {
	resposta, err := http.Get(site)
	if err != nil {
		log.Fatalln("Ocorreu um erro:", err)
	}
	if resposta.StatusCode == http.StatusOK {
		fmt.Println("Site: ", site, " foi carregado corretamente")
		registrarLog(site, true)
	} else {
		fmt.Println("Site: ", site, " está com problemas: ", resposta.StatusCode)
		registrarLog(site, false)
	}

}

func registrarLog(site string, status bool) {
	var statusLinha string
	if status {
		statusLinha = "Online  - "
	} else {
		statusLinha = "Offline - "
	}
	arquivo, err := os.OpenFile("siteStatus.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer arquivo.Close()
	if err != nil {
		log.Fatalln("Ocorreu um erro ao abrir o arquivo de logs:", err)
	} else {
		arquivo.WriteString("[" + time.Now().Format("02/01/2006 15:04:05") + "]: " + statusLinha + site + "\n")
	}
}

func exibirLogs() {
	arquivo, err := ioutil.ReadFile("siteStatus.log")
	//defer arquivo.Close() nesse caso não é necessário pois ioutil.ReadFile devolve apenas os bytes do arquivo e ja fecha o arquivo
	if err != nil {
		log.Fatalln("Ocorreu um erro ao tentar ler o arquivo de logs:", err)
	}
	fmt.Println(string(arquivo))
}

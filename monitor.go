package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 3
const delay = 5 * time.Second
const formdata = "02/01/2006 15:04:05"

func main() {
	introducao()
	for {
		exibeMenu()
		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 0:
			fmt.Println("Saindo do programa...")
			os.Exit(0)
		default:
			fmt.Println("Comando não reconhecido!")
			os.Exit(-1)
		}
	}
}

func introducao() {
	versao := 1.1
	fmt.Println("Monitorador de Sites")
	fmt.Println("Este programa está na versão", versao)
}
func exibeMenu() {
	fmt.Println("1- Iniciar Monitoramanto")
	fmt.Println("2- Exibir Logs")
	fmt.Println("0- Sair do Programa")
}
func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	return comandoLido
}
func trataErro(err error) int {
	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
		return 0
	}
	return -1
}

func iniciarMonitoramento() {
	sites := leSitesArquivo()

	for i := 0; i < monitoramentos; i++ {
		for i, site := range sites {
			fmt.Println("Testando site", i, ":", site)
			testaSite(site)
		}
		fmt.Println("")
		time.Sleep(delay)
	}
	fmt.Println("")
}

func testaSite(site string) {
	resp, err := http.Get(site)
	trataErro(err)

	if resp.StatusCode == 200 {
		fmt.Println("Site:", site, "foi carregado com sucesso!")
		registraLog(site, true)
	} else {
		fmt.Println("Site:", site, "retornou falha no carregamento! Status code:", resp.StatusCode)
		registraLog(site, false)
	}
}

func leSitesArquivo() []string {
	var sites []string

	arquivo, err := os.Open("sites.txt")
	trataErro(err)

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		fmt.Println(linha)

		sites = append(sites, linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()
	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	trataErro(err)

	arquivo.WriteString(time.Now().Format(formdata) + " - " + site + " - online: " + strconv.FormatBool(status) + "\n")

	arquivo.Close()
}

func imprimeLogs() {
	arquivo, err := ioutil.ReadFile("log.txt")
	trataErro(err)

	fmt.Println(string(arquivo))
}

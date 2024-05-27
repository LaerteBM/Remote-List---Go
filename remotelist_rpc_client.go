package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
)

// Estrutura para argumentos de chamadas RPC.
type Args struct {
	NomeDaLista string
	Valor       int
}

// Estrutura que representa uma lista de inteiros.
type List struct {
	Nome  string
	Itens []int
}

// Função para tentar reconectar ao servidor.
func tryReconnect() *rpc.Client {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Deseja tentar reconectar ao servidor? (s/n): ")
		if !scanner.Scan() {
			break
		}
		response := scanner.Text()
		if response == "s" || response == "S" {
			client, err := rpc.Dial("tcp", "127.0.0.1:1234")
			if err != nil {
				fmt.Println("Erro ao tentar reconectar ao servidor:", err)
				continue
			}
			fmt.Println("Reconectado ao servidor com sucesso.")
			return client
		} else if response == "n" || response == "N" {
			fmt.Println("Encerrando o programa.")
			os.Exit(0)
		} else {
			fmt.Println("Resposta inválida. Por favor, responda 's' ou 'n'.")
		}
	}
	return nil
}

// Função para adicionar um valor à lista no servidor.
func appendToList(client *rpc.Client, nome string, valor int) {
	args := Args{NomeDaLista: nome, Valor: valor}
	var reply bool
	err := client.Call("RemoteList.Append", args, &reply)
	if err != nil {
		fmt.Println("Erro ao adicionar valor à lista:", err)
		client = tryReconnect()
		if client != nil {
			appendToList(client, nome, valor)
		}
		return
	}
	if reply {
		fmt.Printf("Valor %d adicionado à lista %s\n", valor, nome)
	} else {
		fmt.Println("Falha ao adicionar valor à lista")
	}
}

// Função para obter o último valor de uma lista no servidor.
func getLastValue(client *rpc.Client, nome string) {
	var reply int
	err := client.Call("RemoteList.Get", nome, &reply)
	if err != nil {
		fmt.Println("Erro ao obter valor da lista:", err)
		client = tryReconnect()
		if client != nil {
			getLastValue(client, nome)
		}
		return
	}
	fmt.Printf("Último valor da lista %s: %d\n", nome, reply)
}

// Função para remover o último valor de uma lista no servidor.
func removeFromList(client *rpc.Client, nome string) {
	var reply int
	err := client.Call("RemoteList.Remove", nome, &reply)
	if err != nil {
		fmt.Println("Erro ao remover valor da lista:", err)
		client = tryReconnect()
		if client != nil {
			removeFromList(client, nome)
		}
		return
	}
	fmt.Printf("Valor %d removido da lista %s\n", reply, nome)
}

// Função para obter o tamanho de uma lista no servidor.
func sizeOfList(client *rpc.Client, nome string) {
	var reply int
	err := client.Call("RemoteList.Size", nome, &reply)
	if err != nil {
		fmt.Println("Erro ao obter tamanho da lista:", err)
		client = tryReconnect()
		if client != nil {
			sizeOfList(client, nome)
		}
		return
	}
	fmt.Printf("Tamanho da lista %s: %d\n", nome, reply)
}

// Função para listar todos os nomes das listas disponíveis no servidor.
func listAllLists(client *rpc.Client) {
	var reply []string
	err := client.Call("RemoteList.ListAllLists", struct{}{}, &reply)
	if err != nil {
		fmt.Println("Erro ao listar todas as listas:", err)
		client = tryReconnect()
		if client != nil {
			listAllLists(client)
		}
		return
	}
	fmt.Println("Listas disponíveis no servidor:")
	for _, nome := range reply {
		fmt.Println(nome)
	}
}

// Função principal do cliente.
func main() {
	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		client = tryReconnect()
		if client == nil {
			return
		}
	}
	defer client.Close()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Menu:")
		fmt.Println("1 - Adicionar valor à lista")
		fmt.Println("2 - Obter último valor da lista")
		fmt.Println("3 - Remover último valor da lista")
		fmt.Println("4 - Obter tamanho da lista")
		fmt.Println("5 - Listar todas as listas")
		fmt.Println("6 - Sair")
		fmt.Print("Escolha uma opção: ")

		if !scanner.Scan() {
			break
		}
		option, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Opção inválida, tente novamente.")
			continue
		}

		switch option {
		case 1:
			fmt.Print("Informe o nome da lista: ")
			if !scanner.Scan() {
				break
			}
			nome := scanner.Text()

			fmt.Print("Informe o valor a ser adicionado: ")
			if !scanner.Scan() {
				break
			}
			valor, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Valor inválido, tente novamente.")
				continue
			}

			appendToList(client, nome, valor)

		case 2:
			fmt.Print("Informe o nome da lista: ")
			if !scanner.Scan() {
				break
			}
			nome := scanner.Text()

			getLastValue(client, nome)

		case 3:
			fmt.Print("Informe o nome da lista: ")
			if !scanner.Scan() {
				break
			}
			nome := scanner.Text()

			removeFromList(client, nome)

		case 4:
			fmt.Print("Informe o nome da lista: ")
			if !scanner.Scan() {
				break
			}
			nome := scanner.Text()

			sizeOfList(client, nome)

		case 5:
			listAllLists(client)

		case 6:
			fmt.Println("Saindo do programa.")
			return

		default:
			fmt.Println("Opção inválida, tente novamente.")
		}
	}
}

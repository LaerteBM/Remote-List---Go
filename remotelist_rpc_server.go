package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"sync"
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

// Estrutura que gerencia listas remotamente.
type RemoteList struct {
	mu sync.Mutex
}

// Função para carregar listas de um arquivo.
func (rl *RemoteList) LoadFromFile() map[string]*List {
	file, err := os.Open("data.json")
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return make(map[string]*List)
	}
	defer file.Close()

	var lists map[string]*List
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&lists)
	if err != nil {
		fmt.Println("Erro ao decodificar o arquivo:", err)
		return make(map[string]*List)
	}
	return lists
}

// Função para salvar listas em um arquivo.
func (rl *RemoteList) SaveToFile(lists map[string]*List) {
	file, err := os.Create("data.json")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(lists)
	if err != nil {
		fmt.Println("Erro ao codificar o arquivo:", err)
	}
}

// Função para adicionar um valor a uma lista.
func (rl *RemoteList) Append(args Args, reply *bool) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lists := rl.LoadFromFile()
	lista, exists := lists[args.NomeDaLista]
	if !exists {
		lista = &List{Nome: args.NomeDaLista}
		lists[args.NomeDaLista] = lista
	}
	lista.Itens = append(lista.Itens, args.Valor)
	rl.SaveToFile(lists)

	fmt.Printf("Valor %d adicionado à lista %s\n", args.Valor, args.NomeDaLista)
	*reply = true
	return nil
}

// Função para obter o último valor de uma lista.
func (rl *RemoteList) Get(nomeDaLista string, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lists := rl.LoadFromFile()
	lista, exists := lists[nomeDaLista]
	if !exists || len(lista.Itens) == 0 {
		return fmt.Errorf("lista com nome %s não existe ou está vazia", nomeDaLista)
	}

	*reply = lista.Itens[len(lista.Itens)-1]
	fmt.Printf("Último valor da lista %s obtido: %d\n", nomeDaLista, *reply)
	return nil
}

// Função para remover o último valor de uma lista.
func (rl *RemoteList) Remove(nomeDaLista string, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lists := rl.LoadFromFile()
	lista, exists := lists[nomeDaLista]
	if !exists || len(lista.Itens) == 0 {
		return fmt.Errorf("lista com nome %s não existe ou está vazia", nomeDaLista)
	}

	*reply = lista.Itens[len(lista.Itens)-1]
	lista.Itens = lista.Itens[:len(lista.Itens)-1]

	fmt.Printf("Último valor da lista %s removido: %d\n", nomeDaLista, *reply)
	rl.SaveToFile(lists)
	return nil
}

// Função para obter o tamanho de uma lista.
func (rl *RemoteList) Size(nomeDaLista string, reply *int) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lists := rl.LoadFromFile()
	lista, exists := lists[nomeDaLista]
	if !exists {
		return fmt.Errorf("lista com nome %s não existe", nomeDaLista)
	}

	*reply = len(lista.Itens)
	fmt.Printf("Tamanho da lista %s obtido: %d\n", nomeDaLista, *reply)
	return nil
}

// Função para listar todos os nomes das listas disponíveis.
func (rl *RemoteList) ListAllLists(_ struct{}, reply *[]string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lists := rl.LoadFromFile()
	var listNames []string
	for nome := range lists {
		listNames = append(listNames, nome)
	}
	*reply = listNames

	fmt.Println("Listas disponíveis no servidor:", listNames)
	return nil
}

func main() {
	rl := &RemoteList{}
	rpc.Register(rl)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor iniciado. Aguardando conexões...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		fmt.Println("Novo dispositivo conectado:", conn.RemoteAddr().String())
		go rpc.ServeConn(conn)
	}
}

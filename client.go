package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "sync"
)


func gestionErreur(err error) {
    if err != nil {
        panic(err)
    }
}

const (
    IP   = "127.0.0.01" // IP local
    PORT = "3569"       // Port utilisé
	CONN_TYPE = "tcp"
)

func main() {

    var wg sync.WaitGroup

    // Connexion au serveur
    conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%s", IP, PORT))
    gestionErreur(err)

    wg.Add(2)

    go func() { // goroutine dédiée à l'entrée utilisateur
        defer wg.Done()
        for {
            reader := bufio.NewReader(os.Stdin)
            text, err := reader.ReadString('\n')
            gestionErreur(err)

            conn.Write([]byte(text))
        }
    }()

    go func() { // goroutine dédiée à la reception des messages du serveur
        defer wg.Done()
        for {
            message, err := bufio.NewReader(conn).ReadString('\n')
            gestionErreur(err)

            fmt.Print("serveur : " + message)
        }
    }()

    wg.Wait()

}

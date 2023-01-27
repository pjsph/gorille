package main

import (
    "bufio"
    "fmt"
    "net"
)

func gestionErreur(err error) {
    if err != nil {
        panic(err)
    }
}

const (
    IP   = "127.0.0.01"
    PORT = "3569"
	CONN_TYPE = "tcp"
)

func read(conn net.Conn) {
    message, err := bufio.NewReader(conn).ReadString('\n')
    gestionErreur(err)

    fmt.Print("Client:", string(message))

}

func main() {

    fmt.Println("Lancement du serveur ...")

    ln, err := net.Listen(CONN_TYPE, fmt.Sprintf("%s:%s", IP, PORT))
    gestionErreur(err)

    var clients []net.Conn // tableau de clients

    for {
        conn, err := ln.Accept()
        if err == nil {
            clients = append(clients, conn) //quand un client se connecte on le rajoute à notre tableau
        }
        gestionErreur(err)
        fmt.Println("Un client est connecté depuis", conn.RemoteAddr())

        go func() { // création de notre goroutine quand un client est connecté
            buf := bufio.NewReader(conn)

            for {
                name, err := buf.ReadString('\n')
				fmt.Print("from ",conn.RemoteAddr()," Client:", string(name))
                if err != nil {
                    fmt.Printf(" disconnected.\n")
                    break
                }
                conn.Write([]byte(name)) // on envoie un message à chaque client
            }
        }()
    }
}

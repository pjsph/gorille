package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "io"
    "strings"
    //"sync"
)


func gestionErreur(err error) {
    if err != nil {
        panic(err)
    }
}

func fillString(input string, maxLength int) string {
    for {
        length := len(input)
        if length < maxLength {
            input += ":"
        } else {
            break
        }
    }
    return input
}

const (
    IP   = "127.0.0.1" // IP local
    PORT = "3569"       // Port utilisé
	CONN_TYPE = "tcp"
    BUFFER_SIZE = 1024
)

func main() {
    conn, err := net.Dial(CONN_TYPE, fmt.Sprintf("%s:%s", IP, PORT))
    gestionErreur(err)
    defer conn.Close()
    fmt.Println("Connecte au serveur. Rentrez le chemin vers le fichier wav a envoyer:")

    var fileName string
    fmt.Scanln(&fileName)

    file, err := os.Open(fileName)
    defer file.Close()
    gestionErreur(err)

    sendToServer(conn, file)

    outFile, err := os.Create(fileName[0:len(fileName) - 4] + "_out.wav")
    defer outFile.Close()
    gestionErreur(err)

    receiveFromServer(conn, outFile)
}

func sendToServer(conn net.Conn, file *os.File) {
    fileInfo, err := file.Stat()
    gestionErreur(err)

    fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)

    fmt.Println("Sending fileSize...")
    conn.Write([]byte(fileSize))

    sendBuffer := make([]byte, BUFFER_SIZE)
    fmt.Println("Sending file...")
    count1, count2 := 0, 0
    for {
        nb, err := file.Read(sendBuffer)
        count1 += nb
        if err == io.EOF {
            break
        }
        nb2, _ := conn.Write(sendBuffer)
        count2 += nb2
    }
    fmt.Println("Read:", count1, "bytes and sent", count2, "bytes")
}

func receiveFromServer(conn net.Conn, outFile *os.File) {
    bufferFileSize := make([]byte, 10)

    conn.Read(bufferFileSize)
    fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

    var receivedBytes int64
    count := 0
    for {
        if (fileSize - receivedBytes) < BUFFER_SIZE {
            nb, _ := io.CopyN(outFile, conn, fileSize - receivedBytes)
            count += int(nb)
            conn.Read(make([]byte, (receivedBytes + BUFFER_SIZE) - fileSize))
            break
        }
        nb, _ := io.CopyN(outFile, conn, BUFFER_SIZE)
        count += int(nb)
        receivedBytes += BUFFER_SIZE
    }

    fmt.Println("Received and written", count, "bytes")
}

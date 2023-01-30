package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
    "bytes"
    "math"
    "time"
    "sync"
    "io"

	pkg "server/GoAudio/wave"
)

func gestionErreur(err error) {
    if err != nil {
        panic(err)
    }
}

const (
    IP   = "127.0.0.1"
    PORT = "3569"
	CONN_TYPE = "tcp"
    BUFFER_SIZE = 1024
)

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

func read(conn net.Conn) {
    message, err := bufio.NewReader(conn).ReadString('\n')
    gestionErreur(err)

    fmt.Print("Client:", string(message))

}

var res []pkg.Frame

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}


func handle(conn net.Conn) {
    wave := receiveFromClient(conn)

    SAMPLES_SIZE := len(wave.Frames) / 16

    var wg sync.WaitGroup
    start_t := time.Now()
    res = make([]pkg.Frame, len(wave.Frames))
    var ires, count int = 0, 0
    for i := 0; i < len(wave.Frames); i += (SAMPLES_SIZE) {
        var to_compute []pkg.Frame = wave.Frames[i:min(i+(SAMPLES_SIZE), len(wave.Frames))]
        wg.Add(1)
        go distort(i, to_compute, 0.1, &wg)
        ires += min(SAMPLES_SIZE, len(wave.Frames)-i)
        count++
    }

    wg.Wait()
    end_t := time.Now()
    fmt.Println("Finished, computed", ires, "samples in", end_t.Sub(start_t), "on", count, "goroutines")
    
    sendToClient(conn, wave, res)
}

func sendToClient(conn net.Conn, wave pkg.Wave, res []pkg.Frame) {
    var buf bytes.Buffer
    pkg.WriteWaveToWriter(res, wave.WaveFmt, &buf)

    fmt.Println("Sending fileSize...")
    newFileSize := fillString(strconv.FormatInt(int64(buf.Len()), 10), 10)
    conn.Write([]byte(newFileSize))

    sendBuffer := make([]byte, BUFFER_SIZE)
    fmt.Println("Sending data...")
    count1, count2 := 0, 0
    for {
        nb1, err := buf.Read(sendBuffer)
        if err == io.EOF {
            break
        }
        nb2, _ := conn.Write(sendBuffer)
        count1 += nb1
        count2 += nb2
    }
    fmt.Println("Read", count1, "bytes and sent", count2, "bytes")
    fmt.Println("done")
}

func receiveFromClient(conn net.Conn) pkg.Wave {
    bufferFileSize := make([]byte, 10)

    conn.Read(bufferFileSize)
    fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

    var receivedBytes int64
    buff := make([]byte, fileSize)
    for {
        if (fileSize - receivedBytes) < BUFFER_SIZE {
            slice := buff[receivedBytes:fileSize]
            conn.Read(slice)
            conn.Read(make([]byte, (receivedBytes + BUFFER_SIZE) - fileSize))
            break
        }
        slice := buff[receivedBytes:(receivedBytes + BUFFER_SIZE)]
        conn.Read(slice)
        receivedBytes += BUFFER_SIZE
    }

    fmt.Println("Received", len(buff), "bytes")

    reader := bytes.NewReader(buff)

    fmt.Println("Parsing wave file..")

    wave, err := pkg.ReadWaveFromReader(reader)
    if err != nil {
        panic("Could not parse wave file")
    }
    
    fmt.Printf("Read %v samples\n", len(wave.Frames))
    return wave
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

        go handle(conn)
    }
}

func changeAmplitude(startIndex int, samples []pkg.Frame, scalefactor float64, wg *sync.WaitGroup) {
    for i, s := range samples {
        res[startIndex + i] = pkg.Frame(float64(s) * scalefactor)
    }
    wg.Done()
}

func distort(startIndex int, samples []pkg.Frame, drive float64, wg *sync.WaitGroup) {
    for i, s := range samples {
        r := float64(s) * math.Pow(10, 2*drive)
        r = math.Max(-1, math.Min(1, r))
        r = r - r*r*r/3
        res[startIndex + i] = pkg.Frame(r * 5)
    }
    wg.Done()
}

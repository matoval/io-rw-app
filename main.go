package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"

	_ "github.com/grafana/pyroscope-go/godeltaprof/http/pprof"
)

var workUnits []string

func main() {
	go writeToSocket("/tmp/io-rw-app.sock")
	listenSocket("/tmp/io-rw-app.sock")
}

func writeToSocket(socketPath string) {
	time.Sleep(2 * time.Second)
	// Connect to the socket
	for {
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			fmt.Println("Error connecting to socket:", err)
			os.Exit(1)
		}
		go sendLongMessage(conn)
	}
}

func sendLongMessage(conn net.Conn){
	// Send a message
	message := randomString(99999)
	time.Sleep(1 * time.Second)
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error writing to socket:", err)
		os.Exit(1)
	}
	fmt.Println("Message sent successfully")
	conn.Close()
}

func listenSocket(socketPath string) {
	os.Remove(socketPath)
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("error creating socket: ", err)
	}
	defer listener.Close()
	fmt.Println("listening on: ", socketPath)
	
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error accepting connection")
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	randStr := randomString(8)
	if _, err := os.Stat("/tmp/io-rw-app/"); os.IsNotExist(err) {
		err := os.Mkdir("/tmp/io-rw-app/", 0755)
		if err != nil {
			fmt.Println("error creating directory: ", err)
		}
	}
	createFile("/tmp/io-rw-app/" + randStr + "/")
	workUnits = append(workUnits, randStr)
	fmt.Println("Created work dir " + randStr)

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("error reading: ", err)
		return
	}
	
	err = os.WriteFile("/tmp/io-rw-app/" + randStr + "/" + "stdin", buf, 0644)
	if err != nil {
		fmt.Println("error opening stdin: ", err)
		return
	}
	
	err = os.WriteFile("/tmp/io-rw-app/" + randStr + "/" + "status", []byte("stdin written successfully\n"), 0644)
	if err != nil {
		fmt.Println("error opening status: ", err)
		return
	}
	r, err := os.ReadFile("/tmp/io-rw-app/" + randStr + "/" + "stdin")
	if err != nil {
		fmt.Println("error reading stdin: ", err)
		return
	}
	err = os.WriteFile("/tmp/io-rw-app/" + randStr + "/" + "stdout", r, 0644)
	if err != nil {
		fmt.Println("error opening stdout: ", err)
		return
	}
}

func createFile(filePath string) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		err := os.Mkdir(filePath, os.ModePerm)
		if err != nil {
			fmt.Println("error creating directory: ", err)
		}
	}
	fileNames := []string{"stdin", "stdout", "status"}
	for _, fileName := range fileNames {
		file, err := os.Create(filepath.Join(filePath + fileName))
		if err!= nil {
			fmt.Println("error creating file: ", err)
		}
		defer file.Close()
	}
}

func randomString(length int) string {
	if length < 0 {
		return ""
	}
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randbytes := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		randbytes = append(randbytes, charset[idx.Int64()])
	}

	return string(randbytes)
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func setupPipes(s *ssh.Session, wrc chan []byte, done chan int) {
	//-------
	//
	//-------

	stdinBuf, err := s.StdinPipe()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	stdoutBuf, err := s.StdoutPipe()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	stderrBuf, err := s.StderrPipe()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// write to stdin of the shell
	go func() {
		for {
			select {
			case d := <-wrc:
				_, err := stdinBuf.Write(d)
				if err != nil {
					log.Fatal(err)
					panic(err)
				}
			}
		}
	}()

	// pipe the stdout of shell to go program stdout and stderr
	go func() {
		scanner := bufio.NewScanner(stdoutBuf)

		for {
			if scanner.Scan() {
				rcv := scanner.Bytes()

				raw := make([]byte, len(rcv))
				copy(raw, rcv)

				log.Printf("stdout:::%s", string(raw))
			} else {
				if scanner.Err() != nil {
					log.Printf(scanner.Err().Error())
				} else {
					log.Printf("io.EOF")
					//os.Exit(0)

				}
				done <- 1
			}
		}
	}()

	// pipe the stderr the shell to go program stdout and stderr
	go func() {
		scanner := bufio.NewScanner(stderrBuf)

		for {
			if scanner.Scan() {
				log.Printf("stderr:::%s", string(scanner.Bytes()))
			} else {
				if scanner.Err() != nil {
					log.Printf(scanner.Err().Error())
				} else {
					log.Print("io.EOF")
					//os.Exit(0)
				}
				done <- 1

			}
		}

	}()
}

func main() {
	// hostKeyCallback, err := knownhosts.New("/home/yakul/.ssh/known_hosts")
	// if err != nil {
	// 	log.Fatal(err)
	// 	panic(err)
	// }

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("mypassword"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	stringFlag := flag.String("containerid", "", "container id which will run the commands")
	flag.Parse()

	// ssh into a remote server for now
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", *stringFlag), config)
	if err != nil {

		log.Fatal(err)
		panic(err)
	}

	defer conn.Close()

	s, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer s.Close()

	wrc := make(chan []byte)
	done := make(chan int, 1)

	setupPipes(s, wrc, done)

	err = s.Start("/bin/bash")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	// put into stdinBuf

	fmt.Println("$")

	file, err := os.Open("./install_curl.sh")

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		wrc <- []byte(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Printf("here")

	<-done
	<-done

}

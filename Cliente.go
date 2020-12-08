package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"strconv"
	"time"
)

type Proceso struct {
	Id int
	Tiempo uint64
	Seguir bool
}

func (proc *Proceso) ejecutarProceso() {
	for proc.Seguir {
		infoProc := strconv.Itoa(proc.Id) + ":" + strconv.Itoa(int(proc.Tiempo))
		fmt.Println(infoProc)
		proc.Tiempo++
		time.Sleep(time.Second / 2)
		
	}
}

func (proc *Proceso) detener() {
	proc.Seguir = false
}

func recibirProceso() {
	var proc Proceso
	var input string

	c, err := net.Dial("tcp", ":9999")	
	if err != nil {
		fmt.Println(err)
		return
	}

	err = gob.NewDecoder(c).Decode(&proc)
	if err != nil {
		fmt.Println(err)
		return
	}

	if proc.Id == -1 {
		fmt.Println("No hay procesos disponibles para recibir")
		return
	}

	proc.Seguir = true
	proc.Tiempo++

	go proc.ejecutarProceso()
	fmt.Scanln(&input)
	proc.Tiempo--

	err = gob.NewEncoder(c).Encode(proc)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Close()
}

func main() {
	recibirProceso()
	fmt.Println("Terminar Cliente")
}


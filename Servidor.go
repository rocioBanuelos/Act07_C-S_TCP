package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"strconv"
	"time"
	"container/list"
)

type AdminProcesos struct {
	procesos list.List
	mostrar  bool
	contId   int
	canal    chan string
}

type Proceso struct {
	Id       int
	Tiempo   uint64
	Seguir   bool
	Regresar bool
	Terminar bool
}

func (adminProc *AdminProcesos) agregarProceso() {
	proc := new(Proceso)
	proc.Tiempo = 0
	proc.Seguir = true
	proc.Id = adminProc.contId
	proc.Terminar = false
	proc.Regresar = false
	
	adminProc.procesos.PushBack(proc)
	adminProc.contId++
}

func (proc *Proceso) detener() {
	proc.Seguir = false
}

func servidor(adminProc *AdminProcesos) {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleCliente(c, adminProc)
	}
}

func handleCliente(c net.Conn, adminProc *AdminProcesos) {
	if adminProc.procesos.Len() > 0 {
		e := adminProc.procesos.Front()
		proc := e.Value.(*Proceso)
		
		err := gob.NewEncoder(c).Encode(proc)
		if err != nil {
			fmt.Println(err)
			return
		}
		proc.Terminar = true
		
		var procNuevo *Proceso
		err = gob.NewDecoder(c).Decode(&procNuevo)
		if err != nil {
			fmt.Println(err)
			return
		}

		procNuevo.Regresar = true
		adminProc.procesos.PushBack(procNuevo)
	} else {
		var procNuevo Proceso
		procNuevo = Proceso{-1, 0, false, false, false}
		
		err := gob.NewEncoder(c).Encode(procNuevo)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func correrProcesos(adminProc *AdminProcesos) {
	for {
		for e := adminProc.procesos.Front(); e != nil; e = e.Next() {
			proc := e.Value.(*Proceso)
			if proc.Seguir {
				infoProc := strconv.Itoa(proc.Id) + ":" + strconv.Itoa(int(proc.Tiempo))
				proc.Tiempo++
				if adminProc.mostrar {
					if proc.Seguir {
						if proc.Regresar {
							proc.Tiempo--
							fmt.Println(infoProc, "Regresar proceso")
							proc.Regresar = false
						} else {
							if proc.Terminar {
								fmt.Println(infoProc, "Salir proceso")
								proc.Terminar = false
								proc.detener()
								e := adminProc.procesos.Front()
								adminProc.procesos.Remove(e)
							} else {
								fmt.Println(infoProc)
							}
						}
					}
				}
			}
		}
		time.Sleep(time.Second / 2)
		fmt.Println("***************")
	}
}

const id_INICIAL int = 0

func main() {
	var input string
	adminProc := new(AdminProcesos)
	adminProc.contId = id_INICIAL
	adminProc.canal = make(chan string)

	adminProc.agregarProceso()
	adminProc.agregarProceso()
	adminProc.agregarProceso()
	adminProc.agregarProceso()
	adminProc.agregarProceso()

	go servidor(adminProc)

	adminProc.mostrar = true
	go correrProcesos(adminProc)

	fmt.Scanln(&input)
	fmt.Println("Fin de Ejecución del Servidor")
}
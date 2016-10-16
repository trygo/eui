package eui

import (
	//"bytes"
	"log"
	"testing"
	//"time"
	//"strings"
	//"unsafe"
)

func Test_Elements(t *testing.T) {
	es := NewElements()
	es.Add(NewElement())
	es.Add(NewElement())
	es.Add(NewElement())
	log.Println("Test_Elements", es.elements)
	es.Add(NewElement(), 1)
	log.Println("Test_Elements", es.elements)
	es.Add(NewElement(), 0)
	log.Println("Test_Elements", es.elements)
	es.Add(NewElement(), -2)
	log.Println("Test_Elements", es.elements)
	e := NewElement()
	es.Add(e, len(es.elements)+1)
	log.Println("Test_Elements", es.elements)
	es.Add(e, len(es.elements))
	log.Println("Test_Elements", es.elements)
	es.Add(e, len(es.elements)+3)
	log.Println("Test_Elements", es.elements)
	es.Add(e, 0)
	log.Println("Test_Elements", es.elements)
	es.Add(e, -3)
	log.Println("Test_Elements", es.elements)
	es.Add(e, 3)
	log.Println("Test_Elements", es.elements)
	es.Add(e, len(es.elements))
	log.Println("Test_Elements", es.elements)

	es.Sort()
	log.Println("Test_Elements.Sort", es.elements)

	//es.Remove(e)
	//log.Println("Test_Elements.Remove", es.elements)

	//log.Printf("Test_Elements.p %p\n", e)
	//log.Println(unsafe.Pointer(e))

}

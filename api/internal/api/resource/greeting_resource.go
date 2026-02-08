package resource

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type GreetingResource struct{}

func NewGreetingResource() *GreetingResource {
	return &GreetingResource{}
}

// Hello handles GET /hello
func (gr *GreetingResource) Hello(w http.ResponseWriter, r *http.Request) {
	gr.printMemory()
	fmt.Fprint(w, "Hello from Go REST")
}

func (gr *GreetingResource) printMemory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// In Java implementation: (total - free) / (1000 * 1000)
	// m.Alloc is the bytes of allocated heap objects.
	usedMB := m.Alloc / (1000 * 1000)
	log.Printf("Meg used=%dM\n", usedMB)
}

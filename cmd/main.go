package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Port to run the node on")
	flag.Parse()

	Init(port)

	http.HandleFunc("/get", GetEndpoint)
	http.HandleFunc("/set", SetEndpoint)
	http.HandleFunc("/delete", DeleteEndpoint)

	fmt.Println(`
╔══════════════════════════════════════════════════════════╗
║                  DISTRIBUTED KV STORE                   ║
║                                                          ║
║    ██╗  ██╗██╗   ██╗    ███████╗████████╗ ██████╗ ██████╗ ██████╗ ███████╗    ║
║    ██║ ██╔╝██║   ██║    ██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗██╔═══╝    ║
║    █████╔╝ ██║   ██║    ███████╗   ██║   ██║   ██║██████╔╝█████╗     ║
║    ██╔═██╗ ╚██╗ ██╔╝    ╚════██║   ██║   ██║   ██║██╔══██╗██╔══╝     ║
║    ██║  ██╗ ╚████╔╝     ███████║   ██║   ╚██████╔╝██║  ██║███████╗   ║
║    ╚═╝  ╚═╝  ╚═══╝      ╚══════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝╚══════╝   ║
║                                                          ║
║               🚀 High Performance Key-Value Storage      ║
║               ⚡ Lightning Fast Distributed Architecture ║
║               🔒 Secure & Reliable Data Management       ║
╚══════════════════════════════════════════════════════════╝`)
	fmt.Println("🌟 Server is running on port " + port)
	fmt.Println("📡 Ready to handle requests...")
	http.ListenAndServe(":"+port, nil)
}

package main

func main() {
	config := readConfig("config.ini")
	server := NewDnsProxyServer(5300, config)
	server.run()
}

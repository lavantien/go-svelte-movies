install:
	go get
	cd webui && npm install && cd ..

server:
	go run main.go

webui:
	cd webui && npm run dev -- --open && cd ..

.PHONY: install server webui

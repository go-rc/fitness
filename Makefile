.PHONY : db

serve:
	@@echo "Listening..."
	@@go run server.go

server:
	go build server.go

importer:
	go build ls_import.go

db:
	cd db && \
	./migrate.sh

dropdb:
	cd db && \
	./migrate.sh -x

all: server importer

clean:
	rm -rf server
	rm -rf ls_import

.PHONY : db

serve:
	@@echo "Listening..."
	@@go run server.go

db:
	cd db && \
	./migrate.sh

dropdb:
	cd db && \
	./migrate.sh -x

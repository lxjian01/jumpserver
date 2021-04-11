VERSION=1.0.0
NAME=jumpserver
ENV=dev

all:clean build package
clean:
	rm -rf distribute
build:
	mkdir -p distribute/${NAME}/config;
	mkdir -p distribute/${NAME}/bin;
	go build -o ${NAME};
	mv ${NAME} distribute/${NAME}/bin;
	cp config/${ENV}/jumpserver.yaml distribute/${NAME}/config/;
	cp scripts/*.sh distribute/${NAME}/bin/;
	chmod +x distribute/${NAME}/bin/*;
package:
	tar czf distribute/${NAME}.tar.gz distribute/${NAME};

.PHONY: all cpu memoria entradasalida kernel

all: cpu memoria entradasalida kernel

cpu:
	cd cpu && mkdir -p bin && go build -o bin/cpu && ./bin/cpu $(ENV)

memoria:
	cd memoria && mkdir -p bin && go build -o bin/memoria && ./bin/memoria $(ENV)

entradasalida:
	cd entradasalida && mkdir -p bin && go build -o bin/entradasalida && ./bin/entradasalida $(ENV)

kernel:
	cd kernel && mkdir -p bin && go build -o bin/kernel && ./bin/kernel $(ENV)

clean:
	rm -rf cpu/bin memoria/bin entradasalida/bin kernel/bin


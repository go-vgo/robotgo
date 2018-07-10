all: callback.go types_auto.go gofmt

install:
	go install -p 6 . ./ewmh ./gopher ./icccm ./keybind ./motif ./mousebind \
		./xcursor ./xevent ./xgraphics ./xinerama ./xprop ./xrect ./xwindow

push:
	git push origin master
	git push github master

build-ex:
	find ./_examples/ -type d -wholename './_examples/[a-z]*' -print0 \
		| xargs -0 go build -p 6

gofmt:
	gofmt -w *.go */*.go _examples/*/*.go
	colcheck *.go */*.go _examples/*/*.go

callback.go:
	scripts/write-events callbacks > xevent/callback.go

types_auto.go:
	scripts/write-events evtypes > xevent/types_auto.go

tags:
	find ./ \( -name '*.go' -and -not -wholename './tests/*' -and -not -wholename './_examples/*' \) -print0 | xargs -0 gotags > TAGS

loc:
	find ./ -name '*.go' -and -not -wholename './tests*' -and -not -name '*keysymdef.go' -and -not -name '*gopher.go' -print | sort | xargs wc -l

ex-%:
	go run _examples/$*/main.go

gopherimg:
	go-bindata -f GopherPng -p gopher -i gopher/gophercolor-small.png -o gopher/gopher.go


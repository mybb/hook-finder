if ! hash gox 2>/dev/null; then
	go get github.com/mitchellh/gox
fi

gox -output="bin/hook-finder-{{.OS}}-{{.Arch}}" -os="windows linux darwin" -arch="386 amd64" github.com/mybb/hook-finder/src
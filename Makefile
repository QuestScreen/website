all: _site

.PHONY: _site generated plugin-tutorial/tutgen

plugin-tutorial/tutgen: plugin-tutorial/tutgen.go
	cd plugin-tutorial && go build -o tutgen

generated: plugin-tutorial/tutgen
	mkdir -p generated
	cd generated && ../plugin-tutorial/tutgen

_site: generated
	bundle exec jekyll build
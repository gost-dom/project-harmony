.PHONY: run
run:
	templ generate --watch --proxy="http://localhost:8081" --cmd="go run ."


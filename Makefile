run-basic-consumer:
	go run basic/consumer/main.go
run-basic-producer:
	go run basic/producer/main.go
run-header-receive: # usage -> make run-header-receive args="-h 'wakwaw'"
	go run header/receive/main.go $(args)
run-header-emit: # usage -> make run-header-emit args="-h 'wakwaw'"
	go run header/emit/main.go $(args)
run-header-receive-example:
	go run header/receive/main.go -h "wakwaw"
run-header-emit-example:
	go run header/emit/main.go -h "wakwaw" -m "bapak mana? bapaka mana? di mana? di jonggol!"
run-pubsub-receive:
	go run pubsub/receive/main.go
run-pubsub-emit:
	go run pubsub/emit/main.go
run-routing-receive: # usage -> make run-routing-receive args="info warning error"
	go run routing/receive_direct/main.go $(args)
run-routing-emit:
	go run routing/emit_direct/main.go
run-routing-receive-example-info:
	go run routing/receive_direct/main.go info
run-routing-receive-example-warning:
	go run routing/receive_direct/main.go warning
run-routing-receive-example-error:
	go run routing/receive_direct/main.go error
run-routing-receive-example-info-warning:
	go run routing/receive_direct/main.go info warning
run-routing-receive-example-info-error:
	go run routing/receive_direct/main.go info error
run-routing-receive-example-warning-error:
	go run routing/receive_direct/main.go warning error
run-routing-receive-example-info-warning-error:
	go run routing/receive_direct/main.go info warning error
run-rpc-server:
	go run rpc/server/main.go
run-rpc-client: # usage -> make run-rpc-client args="'())('"
	go run rpc/client/main.go $(args)
run-rpc-client-example:
	go run rpc/client/main.go "())("
run-topics-receive: # usage -> make run-topics-receive args="lazy.#"
	go run topics/receive/main.go $(args)
run-topics-emit:
	go run topics/emit/main.go
run-topics-receive-example-1:
	go run topics/receive/main.go lazy.#
run-topics-receive-example-2:
	go run topics/receive/main.go *.*.rabbit
run-topics-receive-example-3:
	go run topics/receive/main.go *.orange.*
run-workers-worker:
	go run workers/worker/main.go
run-workers-task:
	go run workers/new_task/main.go
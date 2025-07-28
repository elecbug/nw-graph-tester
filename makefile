.PHONY: venv run graph

venv:
	@if [ ! -d "venv" ]; then \
		echo "Creating Python virtual environment..."; \
		python3 -m venv venv; \
		echo "Installing matplotlib..."; \
		./venv/bin/pip install matplotlib; \
	fi
	@echo "Activating virtual environment..."
	@echo "Run: source venv/bin/activate"

run:
	@echo "Building Go application..."
	go build -o out cmd/main.go
	@echo "Running application..."
	./out

graph:
	@if [ ! -d "venv" ]; then \
		echo "Virtual environment not found. Run 'make venv' first."; \
		exit 1; \
	fi
	@echo "Running graph generator with input=../results/network_metric.jsonl output=../results..."
	cd analyzer && \
		../venv/bin/python dup-p99-graph.py --input ../results/network_metric.jsonl --output ../results/

%:
	@:
all:
	docker build -t goloadrep .
	docker run goloadrep --url=http://google.com --requests=10 --concurrency=10   
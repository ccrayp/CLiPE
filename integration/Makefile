init:
	sudo apt install python3.12-venv -y
	python3 -m venv .venv
	.venv/bin/pip install -r requirements.txt

install: init
	sudo apt install -y libpam0g-dev nlohmann-json3-dev libcurl g++
	sudo cp ca.crt /usr/local/share/ca-certificates/clipe-ca.crt
	sudo update-ca-certificates
	.venv/bin/python scripts/install.py

uninstall: init
	.venv/bin/python scripts/uninstall.py

refresh: init
	.venv/bin/python scripts/refresh.py
init:
	sudo apt install python3.12-venv -y
	python3 -m venv .venv
	.venv/bin/pip install -r requirements.txt

install: init
	.venv/bin/python install.py

uninstall: init
	.venv/bin/python uninstall.py
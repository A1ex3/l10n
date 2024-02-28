init:
	pip install -e ./l10n/
	pip install pytest

test:
	python -m l10n.generator --config="tests/configuration.yml"
	python -m pytest -rA tests
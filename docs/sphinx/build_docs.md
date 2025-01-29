# Build documentation guide

## 1. Prerequisites

```bash
apt install make python3 python3-pip python3-sphinx
```
```bash
python -m pip install sphinx_book_theme myst_parser sphinxcontrib.mermaid sphinx-copybutton
```

## 2. Build documentation (html)

```bash
cd {project_dir}/docs/sphinx
```
```bash
make html
```

## 3. Open built documentation (html)

```bash
cd {project_dir}/docs/_build/html
```

Open index.html via web browser

### 3.1. Alternative run nginx server

```bash
docker run -it --rm -d -p 8080:80 --name web -v ./docs/_build/html:/usr/share/nginx/html nginx
```

Open index.html via web browser using `http://<you-ip-addr>:8080/` or using local address `http://127.0.0.1:8080/`

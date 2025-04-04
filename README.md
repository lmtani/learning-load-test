# Go Load Test

Go Load Test é um aplicativo de linha de comando (CLI) escrito em Go para realizar testes de carga em serviços web. Com ele, é possível definir a URL de destino, o número total de requisições a serem feitas e o nível de concorrência das requisições. Ao término, o aplicativo gera um relatório detalhado dos resultados.

## Recursos

- Definir a URL de destino para o teste.  
- Especificar o número total de requisições.  
- Configurar o nível de concorrência.  
- Gerar um relatório completo com informações da execução.

## Instalação

1. Certifique-se de ter o Go instalado na sua máquina.  

2. Clone o repositório e navegue até o diretório do projeto:

   ```bash
   git clone https://github.com/lmtani/learning-load-test.git
   cd learning-load-test
   ```

3. Compile o aplicativo:

   ```bash
   go build -o loadtest ./cmd/loadtest
   ```

## Uso

Execute o aplicativo diretamente na linha de comando:

```bash
./loadtest --url=<URL> --requests=<TOTAL_DE_REQUISICOES> --concurrency=<NIVEL_DE_CONCORRENCIA>
```

### Exemplo

Para realizar um teste de carga no `http://google.com` com 10 requisições e nível de concorrência de 5:

```bash
./loadtest --url=https://httpbin.org/status/200,404,500 --requests=10 --concurrency=5
```

## Docker

1. Crie a imagem Docker:

   ```bash
   docker build -t go-loadtest .
   ```

2. Execute o contêiner com os parâmetros desejados:

   ```bash
   docker run go-loadtest --url=https://httpbin.org/status/200,404,500 --requests=100 --concurrency=10
   ```

## Relatório

Ao final do teste, o aplicativo exibirá:

- Tempo total de execução.  
- Total de requisições realizadas.  
- Quantidade de requisições com status HTTP 200.  
- Distribuição de outros códigos de status HTTP.  

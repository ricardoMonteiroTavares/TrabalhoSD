# Trabalho Prático 1 - RPC

## Sobre
Disponível em: http://www.ic.uff.br/~simone/sd2/trab1.html

Neste trabalho você irá implementar uma chamada remota a procedimento usando a biblioteca RPyC de Python(https://rpyc.readthedocs.io/en/latest/install.html).

Para instalar a biblioteca rpyc basta executar o comando no terminal:
```
pip install rpyc
```
ou 
```
pip3 install rpyc
```
O grupo deve enviar via **Classroom** um relatório com as respostas às questões até dia **08 de outubro de 2020 às 18:00 horas**.

O modelo que iremos utilizar é orientado a serviços. No arquivo *server.py* apresenta-se um template para um servidor oferecer os serviços de chamada remota de procedimento.

Devem existir os dois métodos para conectar e desconectar uma conexão e os outros métodos podem ser definidos livremente por você. 

Neste caso, você pode escolher os atributos que vão ser expostos para outros processos: se o nome começa com  ***exposed_***, ele poderá ser acessado remotamente, senão ele será acessível somente localmente. 

O programa contido no arquivo *client.py* é um cliente que se conecta com o servidor e executa algumas instruções remotamente. 

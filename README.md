Traba# Trabalho Prático 2 - Raft

## Sobre
Disponível em: http://www.cs.bu.edu/~jappavoo/jappavoo.github.com/451/labs/lab-raft.html

Neste trabalho, você implementará uma parte do Raft, um algoritmo de consenso para gerenciar log replicado.
Para isso, deve ser desenvolvido na linguagem Go e só deve mexer no arquivo *raft.go*, que está localizado em ***src/raft***

Para instalar o compilador Go basta acessar este link: https://golang.org/doc/install

O grupo deve enviar via **Classroom** o código válido até dia **10 de novembro de 2020 às 18:00 horas**.

Implementar eleição de líder e heartbeats (RPCs AppendEntries sem entradas no log) conforme a especificação do Raft. O objetivo é que um único líder seja eleito, que o líder continue sendo o líder se não houver falhas e que um novo líder assuma o controle se o antigo líder falhar ou se os pacotes de/para o antigo líder forem perdidos. Execute *go test -run 2A* para testar seu código.

E para rodar, precisa definir a variável de ambiente GOPATH e o seu valor é a localização deste trabalho antes do *src*
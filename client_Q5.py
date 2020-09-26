import rpyc
import sys
import os
import time

from random import randint

if len(sys.argv) < 2:
    exit("Usage {} SERVER".format(sys.argv[0]))

start = time.time()

server = sys.argv[1]
conn = rpyc.connect(server, 5000)

vector = []
n = 10000

for i in range(n):
    value = randint(0, 1000)
    vector.append(value)

print(conn.root.get_sum_vector(vector))

end = time.time()
print("Tempo de execução no Cliente")
print(str(round(end-start,2))+" seg")
print("\n\n")

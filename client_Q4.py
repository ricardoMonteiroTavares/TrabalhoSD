import rpyc
import sys
import os

if len(sys.argv) < 2:
    exit("Usage {} SERVER".format(sys.argv[0]))

server = sys.argv[1]
conn = rpyc.connect(server, 5000)

vector = []
n = int(input("Enter vector size: "))

for i in range(n):
    value = int(input("Enter value for array: "))
    vector.append(value)

print(conn.root.get_sum_vector(vector))
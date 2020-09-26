# Dockerfile para iniciar servidor "remotamente"

FROM python:alpine3.7

WORKDIR /app

COPY server.py /app

RUN pip install rpyc
EXPOSE 5000
CMD python ./server.py
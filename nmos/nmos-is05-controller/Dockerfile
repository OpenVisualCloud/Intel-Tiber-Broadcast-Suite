FROM python:3.9-slim
WORKDIR /nmos
COPY . /nmos
RUN pip install --upgrade pip && pip install --no-cache-dir -r requirements.txt
ENV RECEIVER_IP=localhost
ENV SENDER_IP=localhost
ENV RECEIVER_PORT=95
ENV SENDER_PORT=90
CMD ["sh", "-c", "python threaded-nmos-controller05.py --receiver_ip $RECEIVER_IP --receiver_port $RECEIVER_PORT --sender_ip $SENDER_IP --sender_port $SENDER_PORT"]

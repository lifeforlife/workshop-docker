FROM python:3.7.5-alpine

# Add non root user
RUN addgroup -S app && adduser app -S -G app
RUN mkdir -p /home/app/function
RUN chown app /home/app

USER app

ENV PATH=$PATH:/home/app/.local/bin

WORKDIR /home/app/

COPY function/* function/
COPY index.py           .
COPY requirements.txt   .

USER root
RUN apk --no-cache add curl \
      make \
      automake \
      gcc \
      python3-dev \
      musl-dev \
      libffi-dev \
    && pip install -r requirements.txt \
    && apk del curl \
      make \
      automake \
      gcc \
      python3-dev \
      musl-dev \
      libffi-dev \
      --no-cache

USER app

EXPOSE 5000

CMD python index.py

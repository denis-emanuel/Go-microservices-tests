#base go image
#build a tiny docker image and just copy the executable

FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp" ]
FROM alpine
ARG TARGET=dev

RUN apk update && \
    apk add tzdata && \
    /bin/cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone

RUN mkdir /app
COPY {{.ProjectName}}  /app/{{.ProjectName}}
COPY conf_$TARGET.yml  /app/conf.yml

EXPOSE {{.Port}}
WORKDIR /app

CMD ["./{{.ProjectName}}"]

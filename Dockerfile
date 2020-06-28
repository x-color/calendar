FROM node:latest as nodebuilder
WORKDIR /app/web/calendar
COPY . /app
RUN npm install
RUN npm run build

FROM golang:latest as gobuilder
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /app
COPY . .
RUN mkdir bin
RUN go build -o bin/calendar .

FROM alpine:latest
WORKDIR /app
COPY --from=gobuilder /app/bin /app/bin
COPY --from=nodebuilder /app/web/calendar/dist /app/web/calendar/dist

CMD /app/bin/calendar
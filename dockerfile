# goのイメージをDockerHubから流用する(Alpine Linux)
FROM golang:1.22.5-alpine3.19
# Linuxパッケージ情報の最新化+gitがないのでgitを入れる
RUN apk update && apk add git
# ログのタイムゾーンを指定
ENV TZ /usr/share/zoneinfo/Asia/Tokyo
# GOPATHを設定し、PATHに追加
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
# コンテナ内の作業ディレクトリを指定
WORKDIR /app
# go.mod と go.sum をコピー
COPY ./app/go.mod ./app/go.sum ./
# ソースコードをコンテナ内にコピー
# COPY ./app .
# /app/go.modに記載された依存関係の解決＋必要なパッケージのダウンロードを実行
RUN go mod download
# Airのバイナリをインストール
RUN go install github.com/cosmtrek/air@v1.52.1
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/99designs/gqlgen@latest
# 依存関係の整理
RUN go mod tidy
# コンテナの公開するポートを指定
EXPOSE 5050
# 起動時のコマンド(airを使用するため)
CMD ["air", "-c", ".air.toml"]
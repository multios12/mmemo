FROM node:22-alpine
RUN apk  --no-cache --update add go git tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    go install -v golang.org/x/tools/gopls@latest && \
    go install -v github.com/go-delve/delve/cmd/dlv@latest && \
    go install -v honnef.co/go/tools/cmd/staticcheck@latest

#ENTRYPOINT [ "yarn", "--cwd", "/workspaces/mmemo/front", "dev" ]
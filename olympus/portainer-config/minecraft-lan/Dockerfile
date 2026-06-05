FROM ubuntu:latest

ENV SERVER=

RUN apt-get update \
	&& apt-get install -y --no-install-recommends wget ca-certificates curl \
	&& rm -rf /var/lib/apt/lists/*

RUN wget -O /usr/local/bin/phantom \
	"$(curl -fsSL https://api.github.com/repos/jhead/phantom/releases \
	| sed -n 's/^[[:space:]]*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' \
	| head -n 1 \
	| sed 's#^#https://github.com/jhead/phantom/releases/download/#; s#$#/phantom-linux#')" \
	&& chmod +x /usr/local/bin/phantom

CMD ["sh", "-c", "exec /usr/local/bin/phantom \"$SERVER\""]
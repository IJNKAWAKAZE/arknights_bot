FROM debian:12
WORKDIR /root
ENV TZ="Asia/Shanghai"

RUN apt update && apt install ffmpeg\
          fonts-noto-*\
          ca-certificates\
          libglib2.0-0\
          libnss3\
          libnspr4\
          libatk1.0-0\
          libatk-bridge2.0-0\
          libcups2\
          libdrm2\
          libdbus-1-3\
          libxcb1\
          libxkbcommon0\
          libx11-6\
          libxcomposite1\
          libxdamage1\
          libxext6\
          libxfixes3\
          libxrandr2\
          libgbm1\
          libpango-1.0-0\
          libcairo2\
          libasound2\
          libatspi2.0-0 -y
COPY arknights /root
CMD ["/root/arknights"]

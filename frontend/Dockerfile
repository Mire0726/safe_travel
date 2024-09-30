# ベースイメージを指定
FROM node:18-alpine

# 作業ディレクトリを設定
WORKDIR /app

# キャッシュ利用で効率化するために別でコピー
COPY package.json yarn.lock ./

# 依存関係をインストール
RUN yarn install --frozen-lockfile

# ソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN yarn build

# ポート番号を指定
EXPOSE 3000

# アプリケーションを起動
CMD ["yarn", "dev"]
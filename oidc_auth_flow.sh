#!/bin/bash

set -e  # エラー時にスクリプトを停止

# デフォルト値
BASE_URL="http://localhost:8080"
CLIENT_ID="my-client"
REDIRECT_URI="http://localhost:8080/callback"
STATE="xyz"
DEFAULT_SCOPE="openid,email"

# 引数から SCOPE を取得（指定がなければデフォルト）
SCOPE="${1:-$DEFAULT_SCOPE}"

echo "🔹 Step 1: 認可コードを取得"
AUTH_URL="${BASE_URL}/authorize?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&state=${STATE}&scope=${SCOPE}"
echo "Request: ${AUTH_URL}"

# 認可コードを取得（リダイレクトURLを解析）
REDIRECT_HEADER=$(curl -s -i -L "${AUTH_URL}")
CODE=$(echo "$REDIRECT_HEADER" | grep -oE "code=[^&]+" | sed 's/code=//')

if [ -z "$CODE" ]; then
    echo "❌ 認可コードの取得に失敗しました。レスポンスヘッダー:"
    echo "$REDIRECT_HEADER"
    exit 1
fi
echo "✅ 認可コード: $CODE"

echo ""
echo "🔹 Step 2: アクセストークンを取得"
TOKEN_RESPONSE=$(curl -s -X POST -d "code=${CODE}" "${BASE_URL}/token")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" == "null" ]; then
    echo "❌ アクセストークンの取得に失敗しました。レスポンス: $TOKEN_RESPONSE"
    exit 1
fi
echo "✅ アクセストークン: $ACCESS_TOKEN"

echo ""
echo "🔹 Step 3: ユーザー情報を取得"
USERINFO_RESPONSE=$(curl -s -H "Authorization: Bearer ${ACCESS_TOKEN}" "${BASE_URL}/userinfo")

echo "✅ ユーザー情報:"
echo "$USERINFO_RESPONSE"

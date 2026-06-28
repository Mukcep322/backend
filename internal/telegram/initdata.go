package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

func ValidateInitData(initData, botToken string) (map[string]string, error) {
	parsed, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("invalid init data: %w", err)
	}

	hash := parsed.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash not found")
	}

	// Проверяем auth_date (не старше 24 часов)
	authDate := parsed.Get("auth_date")
	if authDate != "" {
		var timestamp int64
		fmt.Sscanf(authDate, "%d", &timestamp)
		authTime := time.Unix(timestamp, 0)
		if time.Since(authTime) > 24*time.Hour {
			return nil, fmt.Errorf("init data expired")
		}
	}

	// Удаляем hash из данных
	delete(parsed, "hash")

	// Сортируем параметры
	keys := make([]string, 0, len(parsed))
	for k := range parsed {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Формируем строку для проверки
	var dataCheckString []string
	for _, k := range keys {
		dataCheckString = append(dataCheckString, fmt.Sprintf("%s=%s", k, parsed.Get(k)))
	}
	secretKey := hmacSHA256([]byte(botToken), []byte("WebAppData"))
	computedHash := hmacSHA256(secretKey, []byte(strings.Join(dataCheckString, "\n")))

	if computedHash != hash {
		return nil, fmt.Errorf("invalid hash")
	}

	// Парсим user данные
	userData := make(map[string]string)
	for k, v := range parsed {
		userData[k] = v
	}

	return userData, nil
}

func hmacSHA256(key, data []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

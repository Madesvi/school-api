package middlewares

import (
	"fmt"
	"net/http"
)

// Функция возвращает http.Handler значит исполняет его

func SecurityHeaders(next http.Handler) http.Handler {
	fmt.Println("Security Header Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Security Header Middleware being returned...")
		w.Header().Set("X-DNS-Prefetch-Control", "off")     // Браузер не будет резолвить ссылки заране
		w.Header().Set("X-Frame-Options", "DENY")           // Запрещает встраивание страницы во фреймы clickjacking
		w.Header().Set("X-XSS-Protection", "1; mode=block") // Block XSS atack
		w.Header().Set("X-Content-Type-Options", "nosniff") // MIME тип ресурса запрещает угадывать тип ресурса
		w.Header().Set("Strict Transport Security", "max-age=63072000; includeSubDomains; preload")
		// Заставляет обращаться только по HTTPS
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		// Политика безопасности контента, разрешает загрузку контента только с того же источника
		w.Header().Set("Referrer-Policy", "no-referrer")
		// Запрещает передачу Referer при переходах повышая конфиденциальность
		w.Header().Set("X-Powered-By", "Django")
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")
		next.ServeHTTP(w, r)
		fmt.Println("Security Header ends...")
	})
}

// * Feature-Policy / Permissions-Policy — управляет доступом к функциям браузера (камера, микрофон и т.д.).
// * Cache-Control — управляет кешированием, особенно важно для конфиденциальных страниц (например, no-store).
// * Expect-CT — требует проверки сертификата через Certificate Transparency (устарел, но может использоваться).
// * Cross-Origin-* (например, Cross-Origin-Embedder-Policy, Cross-Origin-Opener-Policy) — защищают от атак типа Spectre.
// * Clear-Site-Data — очищает данные браузера при выходе.
// * Server — рекомендуется убирать или минимизировать информацию о сервере.

// BASIC MIDDLEWARE SCELETON

// func securityHeaders(next http.Handler) http.Handler{
// return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request){
// next.ServeHTTP(w, r)
// 	})
// }

# Repositori: agroconnect-api-gateway
## Pintu Masuk API, Reverse Proxy, JWT Auth Middleware & CORS
**Skema Sertifikasi**: BNSP Full-Stack Developer | **Kandidat**: Fredli Fourqoni

---

### Deskripsi
Layanan API Gateway dibangun menggunakan bahasa **Golang** yang bertindak sebagai *single point of entry* untuk semua lalu lintas HTTP dari Klien Web dan Klien Mobile.

### Fitur Utama:
1. **Reverse Proxy Dispatcher**: Mengarahkan rute permintaan ke microservice tujuan:
   - `/api/products*` -> `agroconnect-service-catalog:8081`
   - `/api/auth*` & `/api/orders*` -> `agroconnect-service-order:8082`
   - `/api/weather*` -> `agroconnect-service-weather:8083`
2. **Keamanan & Otorisasi**: Middleware otentikasi JWT HMAC-SHA256 untuk memvalidasi Bearer Token pada endpoint berproteksi.
3. **CORS & Rate Limiter**: Memastikan kebijakan lintas domain aman bagi peramban modern.
4. **Health Check**: Endpoint `/health` untuk monitoring ketersediaan kontainer.

### Port
- Port Eksternal/Internal: `8080`

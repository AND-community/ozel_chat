# and-plugin-konu-ac — Konu Açma Eklentisi

Yeni forum konusu oluştur, taslak kaydet, kategoriler arasında gezin.

---

## Genel Bakış

`and-plugin-konu-ac` forum yazma işlemlerini yönetir. Ana forum görünümü salt okunurdur; tüm konu oluşturma bu eklenti üzerinden gerçekleşir.

Bu eklenti menüde **görünmez** (`Label` alanı boştur). Forumdan `n` tuşuna basıldığında AND tarafından **doğrudan ana süreç içinde** açılır — ayrı bir binary başlatılmaz, geçiş anında gerçekleşir.

---

## Kurulum

```bash
go build -o and-plugin-konu-ac ./Eklentiler/konu_ac

# Windows
go build -o and-plugin-konu-ac.exe ./Eklentiler/konu_ac
```

Binary AND dizininde bulunmalıdır; AND başlangıçta manifest JSON dosyasını okur.

---

## Nasıl Açılır

Forumu açtıktan sonra `n` tuşuna bas.  
AND seçili kategoriyi otomatik iletir; kategori seçim ekranı atlanır.

---

## Kullanım

### Konu formu

Kategori satırında `◀` / `▶` tuşları ile kategori değiştirilebilir.

| Tuş | İşlev |
|-----|-------|
| `tab` | Başlık ↔ İçerik arasında geç |
| `enter` (Başlık) | İçerik alanına geç |
| `enter` (İçerik) | Alt satıra geç |
| `ctrl+s` | Konuyu gönder |
| `ctrl+p` | Kalıcılık talebini aç/kapat (★ gösterilir) |
| `ctrl+t` | Taslak listesini aç (taslak varsa) |
| `◀` / `▶` | Kategori değiştir |
| `esc` | İçeriği taslak olarak kaydet ve foruma dön |
| `ctrl+c` | Kaydetmeden çık |

### Taslak listesi

| Tuş | İşlev |
|-----|-------|
| `↑` / `↓` ya da `j` / `k` | Taslak seç |
| `enter` | Seçili taslağı forma yükle |
| `d` | Seçili taslağı sil |
| `esc` | Forma dön |

---

## Taslak sistemi

`esc` tuşuna basıldığında başlık veya içerik doluysa taslak otomatik kaydedilir.  
Aynı kategori tekrar açıldığında taslaklar `ctrl+t` ile geri yüklenebilir.

Taslaklar `AND_DATA_DIR/taslaklar_<kategori>.json` dosyalarında saklanır.  
Bu dosyalar yereldir, ağa gönderilmez ve `.gitignore`'da yer alır.

---

## Karakter Sınırları

| Alan | Maksimum |
|------|---------|
| Başlık | 100 karakter |
| İçerik | 2000 karakter |

---

## Kategoriler

| | | | |
|--|--|--|--|
| Python | C / C++ | Rust | Go |
| JavaScript | Java / Kotlin | Yazılım | Web |
| Mobil | Yapay Zeka | Veritabanı | DevOps |
| Linux | Bilişim | Siber Güvenlik | Donanım |
| Oyun Geliştirme | Açık Kaynak | Kariyer | Genel |

---

## Kalıcılık talebi

`ctrl+p` ile aktif edilen kalıcılık talebi, moderatörden TTL muafiyeti isteğidir.  
Form başlığında `★ Kalıcılık talep ediliyor` gösterilir.

Bu yalnızca bir istek olarak iletilir; karar her zaman moderatördedir.

---

## Moderasyon notu

Gönderilen konu hemen ağda yayılır, ancak varsayılan olarak moderasyon kuyruğuna girer ve 5 günlük TTL'e tabidir. Kurucu veya sertifikalı bir moderatör onaylayana kadar konu geçici olarak işaretlenir.

---

## Manifest

```json
{
  "name":        "konu_ac",
  "label":       "",
  "version":     "2.1.0",
  "description": "Forum'da yeni konu oluşturma ve taslak yönetimi (menüde gizli, forumdan n ile açılır)",
  "author":      "AND"
}
```

`label` boş olduğu için bu eklenti ana menüde görünmez.

---

## Kaynak

Kaynak kod: [Eklentiler/konu_ac/main.go](main.go)

---

## Sürüm Geçmişi

| Sürüm | Değişiklik |
|-------|------------|
| 2.1.0 | AND ana sürecinde inline çalışma (ayrı binary başlatılmaz); Tab ile alan geçişi; kategori ◀/▶; Enter ile İçerik'te alt satır |
| 2.0.0 | Bağımsız binary; HTTP IPC ile AND ile iletişim; taslak sistemi |
| 1.0.0 | İlk sürüm |

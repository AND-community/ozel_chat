# and-plugin-ozel-chat — Özel Mesajlaşma

Başka bir AND kullanıcısıyla doğrudan ve özel mesajlaş.

---

## Genel Bakış

`and-plugin-ozel-chat` binary'si, libp2p stream protokolü (`/and/dm/1.0.0`) üzerinden iki peer arasında doğrudan mesajlaşma sağlar.

Mesajlar GossipSub üzerinden yayılmaz; yalnızca hedef peer'a iletilir.  
AND, `/and/dm/1.0.0` akışını dinler ve gelen mesajları eklentiye long-poll HTTP üzerinden iletir.

---

## Kurulum

```bash
go build -o and-plugin-ozel-chat ./Eklentiler/ozel_chat

# Windows
go build -o and-plugin-ozel-chat.exe ./Eklentiler/ozel_chat
```

---

## Kullanım

Ana menüden **Özel Chat**'i seç.

### Adımlar

1. **Peer ID girişi** — hedef kişinin libp2p Peer ID'sini gir (ör. `12D3KooWAbc…`)
2. `enter` ile mesajlaşma ekranına geç
3. Mesajını yaz, `enter` ile gönder

Karşı tarafın Peer ID'si AND ana menüsünde kimlik bilgisi ekranında görünür.

### Peer ID girişi ekranı

| Tuş | İşlev |
|-----|-------|
| Yazma | Peer ID'yi gir |
| `enter` | Bağlantıyı onayla, mesajlaşma ekranına geç |
| `esc` ya da `q` | Ana menüye dön |

### Mesajlaşma ekranı

| Tuş | İşlev |
|-----|-------|
| Yazma | Mesajı yaz |
| `enter` | Mesajı gönder |
| `esc` | Peer ID girişine dön |
| `ctrl+c` | Çık |

---

## Teknik Detaylar

### Protokol

Mesajlar `/and/dm/1.0.0` protokolüyle libp2p stream üzerinden iletilir.  
Her mesaj tek JSON paketi olarak gönderilir:

```json
{
  "from": "kullanici_adi",
  "text": "mesaj metni",
  "at":   "2026-01-01T12:00:00Z"
}
```

### Gelen mesajlar — long poll

Eklenti, gelen mesajları doğrudan dinlemez. Bunun yerine AND'ın DM broker'ı (`dmmgr`) `/and/dm/1.0.0` akışını kaydeder ve gelen mesajları yakalar.

Eklenti her 5 saniyede `GET /api/v1/dm/poll` isteği gönderir.  
AND mesaj gelene kadar bekler; gelince JSON array döner, yoksa 5 saniye sonra boş array döner.

### Güvenlik sınırları

- Gelen mesajlar 16 KB ile sınırlıdır; daha büyük paketler atılır
- Her akış için 30 saniyelik okuma zaman aşımı uygulanır
- Gönderme işlemi için 30 saniyelik zaman aşımı uygulanır

### Veri saklama

Mesaj geçmişi yalnızca oturum boyunca bellekte tutulur.  
AND kapatıldığında veya eklenti çıktığında mesajlar silinir; kalıcı mesajlaşma yoktur.

---

## Sınırlamalar

- **Uçtan uca şifreleme yoktur.** Mesajlar libp2p Noise protokolü ile aktarım katmanında şifrelenir, ancak uygulama düzeyinde ek şifreleme uygulanmaz.
- **Kalıcı iletim yoktur.** Alıcı çevrimdışıysa mesaj iletilmez ve kaybolur.
- **Görünen ad doğrulanmaz.** `From` alanı gönderenin beyanıdır; gerçek kimlik için Peer ID'yi kullan.

---

## Manifest

```json
{
  "name":        "ozel_chat",
  "label":       "Özel Chat",
  "version":     "2.0.0",
  "description": "Peer ID ile doğrudan özel mesajlaşma — libp2p stream (/and/dm/1.0.0)",
  "author":      "AND"
}
```

---

## Kaynak

Kaynak kod: [Eklentiler/ozel_chat/main.go](main.go)

---

## Sürüm Geçmişi

| Sürüm | Değişiklik |
|-------|------------|
| 2.0.0 | Bağımsız binary; long-poll HTTP IPC; dmmgr üzerinden akış proxy'si |
| 1.0.0 | İlk sürüm (gömülü plugin sistemi) |

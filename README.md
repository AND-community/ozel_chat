# ozel_chat — archived

**Özel Chat (private chat) has moved into the
[AND repository](https://github.com/AND-community/AND/tree/main/Eklentiler/ozel_chat).** This
repository holds an old version (July 2026); it is archived and no longer updated.

The current Özel Chat — a conversation list with unread counts, a message history that AND keeps
encrypted on your computer, messages to offline peers that are delivered when they connect, sending
files and approving incoming ones — is built together with AND:

```bash
git clone https://github.com/AND-community/AND.git
cd AND
go build -o Eklentiler/ozel_chat/and-plugin-ozel_chat ./Eklentiler/ozel_chat
```

On Windows add `.exe` to the binary name.

---

## Türkçe

**Özel Chat [AND deposuna](https://github.com/AND-community/AND/tree/main/Eklentiler/ozel_chat)
taşındı.** Bu depo eski bir sürümü (Temmuz 2026) içerir; arşivlenmiştir ve artık güncellenmez.
Güncel Özel Chat — okunmamış sayılı konuşma listesi, AND'ın bilgisayarında şifreli tuttuğu mesaj
geçmişi, çevrimdışı eşe bağlanınca iletilen mesajlar, dosya gönderme ve gelen dosyayı onaylama —
AND ile birlikte, yukarıdaki komutlarla derlenir.

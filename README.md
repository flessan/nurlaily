
<div align="center">

# **NurLaily — Your Daily Draft**

> **Catat jurnal harian lewat terminal, ubah jadi website statis bertema cartoon pop secara instan.**

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-3B82F6?style=flat)

</div>

---

## 🌟 Fitur

- 📝 **Tulis jurnal harian** lewat terminal dengan timestamp otomatis
- 😊 Pelacak **suasana hati** per catatan (`--mood`)
- 🏷 Sistem **tag/kategori** untuk mengelompokkan catatan
- 📋 Lihat daftar semua hari penulisan
- 📖 Baca ulang catatan hari ini langsung di terminal
- 📊 Statistik: total kata, streak, top tags
- 🗑 Hapus catatan per-index atau seluruh hari
- 🌐 Bangun website statis sekali perintah (`build`)
- 👀 Preview langsung di browser (`serve`)
- 🔍 Pencarian & filter tag di website
- 📱 Responsif, mendukung `prefers-reduced-motion`

---

## 🚀 Cara Pakai

### Instalasi

```bash
git clone https://github.com/flessan/nurlaily.git
cd nurlaily
go mod tidy
go build -o laily .
````

### Inisialisasi

```bash
./laily init
```

Membuat folder `drafts/` dan menambahkan catatan contoh untuk hari ini. Jalankan sekali saat pertama kali pakai.

### Menulis Catatan

```bash
# Catatan biasa
./laily draft "Hari ini aku belajar Go dan rasanya menyenangkan."
./laily draft "Makan siang di warung baru dekat kos, nasi gorengnya juara."

# Dengan suasana hati
./laily draft --mood "😊" "Akhirnya bugnya ketemu juga setelah 3 jam"

# Dengan tag
./laily draft --tag go,belajar "Belajar goroutine dan channel"
./laily draft --tag kuliah --mood "😴" "Kelas Algoritma jam 8, ngantuk berat"

# Mood + tag sekaligus
./laily draft --mood "🔥" --tag side-project "Mulai bikin CLI jurnal sendiri"
```

### Melihat Catatan

```bash
# Daftar semua hari yang ada catatannya
./laily list

# Baca catatan hari ini di terminal
./laily today
```

---

## 📊 Statistik

```bash
./laily stats
```

Contoh output:

```
  📅 12 hari menulis
  📝 34 catatan
  📊 2840 kata
  🔥 5 hari (streak saat ini)
  🏆 8 hari (streak terpanjang)
  🏷 #belajar(12), #go(8), #kuliah(6), #side-project(4)
```

---

## 🗑 Menghapus Catatan

```bash
# Hapus catatan ke-0 (indeks mulai dari 0) pada tanggal tertentu
./laily delete 2026-04-23 0

# Hapus seluruh catatan pada tanggal tertentu
./laily delete 2026-04-23
```

---

## 🌍 Membangun Website

```bash
./laily build
```

Perintah ini akan:

1. Membaca semua file `.md` di folder `drafts/`
2. Mengonversinya ke HTML menggunakan Goldmark
3. Menyuntikkan hasilnya ke template bertema **blue cartoon pop**
4. Mendukung pencarian teks dan filter tag di sidebar
5. Menghasilkan `dist/index.html` yang siap di-deploy

Output ke folder lain:

```bash
./laily build -o public/
```

---

## 🌐 Preview di Browser

```bash
./laily serve
```

Bangun website lalu buka `http://localhost:3000` di browser. Tekan `Ctrl+C` untuk berhenti.

---

## 📦 Deploy ke Cloudflare Pages / GitHub Pages

Folder `dist/` adalah output statis murni (satu file `index.html`). Upload ke mana saja:

```bash
# Cloudflare Pages
npx wrangler pages deploy dist/

# Atau hubungkan repo GitHub ke Cloudflare Pages,
# set build command: go build -o laily . && ./laily build
# set output directory: dist
```

---

## 📂 Struktur Proyek

```
nurlaily/
├── main.go                          # Entry point
├── cmd/
│   ├── root.go                      # CLI: draft, list, today, stats, delete, init
│   └── build.go                     # CLI: build, serve
├── internal/
│   ├── model/
│   │   └── model.go                 # Struktur data (Entry, Day, Page, Tag)
│   ├── draft/
│   │   └── draft.go                 # Logika tulis, baca, hapus, parse, stats
│   ├── build/
│   │   └── build.go                 # Konversi md → html (Goldmark) + render
│   └── template/
│       ├── template.go              # Go embed + safe HTML render
│       └── index.html               # Template website (blue cartoon pop)
├── drafts/                          # (auto-generated) File markdown jurnal
├── dist/                            # (auto-generated) Output website statis
├── go.mod
└── README.md
```

---

## 📄 Format File Draft

File markdown di `drafts/` menggunakan format sederhana:

```markdown
# 2026-04-23

## 14:30 😊
Hari ini aku belajar Go dan rasanya menyenangkan. #belajar #go

## 12:15
Makan siang di warung baru dekat kos, nasi gorengnya juara. #makan

## 08:00 😴
Kelas jam 8 pagi, ngantuk berat. #kuliah
```

---

## 🎨 Desain Website

Website output menggunakan tema **blue cartoon pop**:

* Latar **biru muda** (`#F0F7FF`) dengan radial gradient halus
* Kartu putih solid dengan **border tebal 3px** dan **cartoon offset shadow**
* Aksen **biru** (`#3B82F6` → `#1E3A8A`) + **kuning** (`#FACC15`) + **pink** (`#F9A8D4`)
* Tipografi **Fredoka** (bulat, playful) + **Caveat** (handwritten) + **JetBrains Mono**
* Dekorasi mengambang: lingkaran, bintang, dots dengan animasi wiggle/pop/spin
* Hover efek kartu: terangkat + scale + gradient stripe muncul di atas
* Sidebar dengan dashed border, tag pills, dan navigasi hari yang interaktif
* Pencarian teks real-time dan filter tag di sidebar
* IntersectionObserver untuk staggered card animation saat scroll
* Sidebar highlight otomatis berdasarkan posisi scroll
* Format tanggal Indonesia: *Rabu, 23 April 2026*
* Responsif untuk mobile (hamburger menu + overlay)
* Mendukung `prefers-reduced-motion`

---

## 🔖 Referensi Perintah

```
laily init                           Inisialisasi folder drafts dengan contoh
laily draft "pesan"                  Tulis catatan baru
laily draft --mood "😊" "pesan"      Tulis dengan suasana hati
laily draft --tag a,b "pesan"        Tulis dengan tag
laily list                           Daftar semua hari penulisan
laily today                          Baca catatan hari ini
laily stats                          Lihat statistik jurnal
laily delete <tanggal> [indeks]      Hapus catatan
laily build [-o dir]                 Bangun website statis
laily serve [-o dir]                 Preview website di localhost:3000
```

---

## 📝 Lisensi

[MIT](LICENSE)

# NurLaily — Daily Draft

> Catat jurnal harian lewat terminal, ubah jadi website statis yang estetis secara instan.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-22d3ee?style=flat)

## Cara Pakai

### Instalasi

```bash
git clone https://github.com/thio/nurlaily.git
cd nurlaily
go mod tidy
go build -o laily .
```

### Menulis Catatan

```bash
./laily draft "Hari ini aku belajar Go dan rasanya menyenangkan."
./laily draft "Makan siang di warung baru dekat kos, nasi gorengnya juara."
```

Setiap perintah `draft` akan menambahkan catatan ke file markdown hari ini (`drafts/2026-04-23.md`), lengkap dengan timestamp jam:menit.

### Melihat Daftar Draft

```bash
./laily list
```

### Membangun Website

```bash
./laily build
```

Perintah ini akan:
1. Membaca semua file `.md` di folder `drafts/`
2. Mengonversinya ke HTML menggunakan Goldmark
3. Menyuntikkan hasilnya ke template estetik (dark mode + glassmorphism)
4. Menghasilkan `dist/index.html` yang siap di-deploy

### Deploy ke Cloudflare Pages / GitHub Pages

Folder `dist/` adalah output statis murni. Upload isinya ke mana saja:

```bash
# Cloudflare Pages
npx wrangler pages deploy dist/

# Atau hubungkan repo GitHub ke Cloudflare Pages,
# set build command: go build -o laily . && ./laily build
# set output directory: dist
```

## Struktur Proyek

```
nurlaily/
├── main.go                      # Entry point
├── cmd/root.go                  # Definisi CLI (Cobra)
├── internal/
│   ├── draft/draft.go           # Logika tulis & list draft
│   ├── build/build.go           # Logika konversi md → html
│   └── template/
│       ├── template.go          # Go embed + render engine
│       └── index.html           # Template HTML estetik
├── drafts/                      # (auto-generated) File markdown jurnal
├── dist/                        # (auto-generated) Output website statis
├── go.mod
└── README.md
```

## Desain

Website output menggunakan:
- **Dark mode** dengan latar hitam pekat (#050510)
- **Glassmorphism** pada sidebar dan kartu konten
- Aksen **cyan** (#22d3ee) dan **ungu muda** (#c084fc)
- **Ambient orbs** yang bergerak perlahan di background
- Tipografi **Space Grotesk** + **JetBrains Mono**
- Responsif untuk mobile
- Mendukung `prefers-reduced-motion`

## Lisensi

[MIT](LICENSE)
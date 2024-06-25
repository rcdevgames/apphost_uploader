# AppHost Uploader

Halo gaes! Ini dia AppHost Uploader, aplikasi buat upload APK ke server yang simpel dan cepat pake Go. 

## Syarat-syarat

Sebelum mulai, pastiin lu udah punya ini:

- [Go](https://golang.org/doc/install) versi terbaru.

## Cara Pakai

Ikutin langkah-langkah berikut buat jalanin aplikasinya:

### 1. Clone Repo

Clone repo ini ke laptop/PC kamu:

```sh
git clone https://github.com/rcdevgames/apphost_uploader.git
cd apphost_uploader
```

### 2. Setup Dependensi

Jalanin perintah ini buat setup dependensi Go:

```sh
go mod tidy
```

### 3. Jalanin Aplikasinya

Buat jalanin aplikasinya, pake perintah ini:

```sh
go run main.go <folder_apk_lau_bos> <versi>
```

Contohnya nih:

```sh
go run main.go ./apk_gweh_nih_bwang.apk 1.0.0
```

## Kontribusi

Kalo lo mau kontribusi, boleh banget! Fork aja repo ini, bikin branch baru buat fitur atau perbaikan lo, terus kirim pull request ke kita.

## Lisensi

Proyek ini dilisensi di bawah [MIT License](LICENSE).

## Kontak

Ada pertanyaan atau saran? Langsung aja email kita di [rcdev.games@gmail.com](mailto:rcdev.games@gmail.com).

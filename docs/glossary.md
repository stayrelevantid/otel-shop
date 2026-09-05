# Glosarium — OTel-Shop Lab

Kamus istilah untuk pemula. Setiap istilah dijelaskan **seperti menjelaskan ke
anak kecil**, lengkap dengan analogi sehari-hari dan di mana istilah itu
muncul di repo ini. Bacanya tidak perlu berurutan — pakai sebagai kamus.

Daftar isi:
1. [Tracing & OpenTelemetry](#tracing--opentelemetry)
2. [Kubernetes & Container](#kubernetes--container)
3. [Go & Testing](#go--testing)
4. [Tools & Quality](#tools--quality)

---

## Tracing & OpenTelemetry

**OpenTelemetry (OTel)**
Kumpulan alat dan aturan standar supaya semua aplikasi bisa "menceritakan apa
yang mereka lakukan" dengan cara yang sama.
*Analogi:* seperti bahasa Indonesia baku — semua provinsi pakai, semua orang
paham. *Di repo:* package `pkg/telemetry`.

**Distributed tracing**
Cara melihat perjalanan satu permintaan yang melompat-lompat antar banyak
aplikasi.
*Analogi:* paket yang dikirim dari Jakarta ke Solo melewati beberapa kantor
pos — tracing mencatat semua perhentian itu dalam satu riwayat utuh.
*Di repo:* satu checkout melibatkan order → inventory → payment.

**Trace**
Riwayat lengkap satu permintaan dari awal sampai selesai.
*Analogi:* satu buku yang menceritakan satu perjalanan dari depan ke belakang.
*Di repo:* punya satu `traceID`, contoh `d87f3b8f...`.

**Span**
Catatan tentang *satu langkah* di dalam perjalanan.
*Analogi:* satu bab dalam buku perjalanan — "berangkat dari rumah (3 menit)".
*Di repo:* `POST /checkout`, `check-inventory`, `sql.conn.query` semua adalah
span.

**Trace ID**
Nomor unik yang menempel di semua span dari satu permintaan.
*Analogi:* nomor resi — semua kantor pos yang memegang paketmu mencatat nomor
yang sama. *Di repo:* inilah yang membuat span order, inventory, dan payment
bergabung jadi satu trace.

**Span parent–child (anak-pinak)**
Span bisa punya orang tua — langkah kecil ada di dalam langkah besar.
*Analogi:* bab 3 bisa punya sub-bab 3.1, 3.2. *Di repo:* `sql.conn.query`
anak dari `GET /inventory/{id}`, yang anak dari `POST /checkout`.

**Instrumentation**
Proses memasang "kamera" di dalam kode supaya aktivitas tercatat.
*Analogi:* memasang CCTV di tiap ruangan rumah. *Di repo:* otelhttp dan
otelsql = CCTV otomatis; span `validate-order` = kamera yang dipasang manual.

**Instrumentasi otomatis vs manual**
Otomatis = alat yang sekali dipasang langsung merekam semua (middleware).
Manual = kita sendiri yang memilih titik penting untuk direkam.
*Analogi:* CCTV otomatis di pintu vs catatan tangan di buku harian.
*Di repo:* otelhttp/otelsql = otomatis; `validate-order` dst = manual.

**SDK**
Perangkat lunak siap pakai untuk memasang kamera-kamera itu.
*Analogi:* paket CCTV lengkap berisi kamera + perekam + kabel.
*Di repo:* `go.opentelemetry.io/otel/sdk`.

**TracerProvider**
Manajer semua kamera: yang mengatur perekaman dan pengiriman hasilnya.
*Analogi:* pusat monitoring CCTV gedung. *Di repo:* dibuat di
`pkg/telemetry.Init()`, dimatikan lewat `shutdown()`.

**Exporter**
Tukang yang mengantar rekaman dari aplikasi ke tempat penyimpanan.
*Analogi:* kurir yang mengantar video CCTV ke ruang arsip.
*Di repo:* OTLP gRPC exporter → ke `otel-collector:14317`.

**OTLP**
Bahasa/kunci standar untuk mengantar rekaman tersebut.
*Analogi:* format bungkusan standar yang dipahami semua kantor pos.
*Di repo:* `14317` = pintu gRPC, `14318` = pintu HTTP di collector.

**Collector**
Perantara yang menerima rekaman dari banyak aplikasi lalu meneruskannya.
*Analogi:* kantor pos pusat: semua kurir drop-off di sini, diurutkan, lalu
dikirim ke gudang akhir. *Di repo:* deployment `otel-collector`, pipeline:
terima → batch → kirim ke Jaeger.

**Pipeline (receiver → processor → exporter)**
Jalur kerja di dalam collector: pintu masuk → pengolahan → pintu keluar.
*Analogi:* pabrik: bahan masuk, dirakit, produk keluar.
*Di repo:* config collector di `deploy/otel-collector/configmap.yaml`.

**Batch processor**
Menyimpan rekaman sebentar lalu mengirim banyak sekaligus.
*Analogi:* antar jemput anak sekolah — tidak antar satu-satu, tapi kumpulkan
dulu, baru berangkat bareng. *Di repo:* `processors: [batch]`.

**Sampling**
Sengaja hanya menyimpan sebagian rekaman (misal 1 dari 10) supaya tidak
berlebihan.
*Analogi:* CCTV dengan storage terbatas — rekam 1 dari 10 lewatan.
*Di repo:* `OTEL_TRACES_SAMPLER_ARG=0.1` → hanya ~10% checkout yang
menghasilkan trace (bukti Day 7: 3 dari 20).

**ParentBased (sampler)**
Aturan: kalau permintaan datang dengan keputusan "direkam/tidak", ikuti
keputusan itu; keputusan baru hanya untuk permintaan paling awal.
*Analogi:* anak ikut keputusan orang tua — satu keluarga satu pilihan.
*Di repo:* `sdktrace.ParentBased(TraceIDRatioBased(...))` di `pkg/telemetry`.

**Propagator (W3C TraceContext)**
Aturan cara "menempelkan nomor resi" ke header HTTP saat request pindah
service.
*Analogi:* kartu barang bawaan kurir yang wajib dicopy tiap kantor pos.
*Di repo:* header `traceparent` otomatis oleh otelhttp transport.

**Baggage**
Tas kecil berisi data kecil yang dibawa-bawa request saat pindah service —
isinya bisa dibaca semua service.
*Analogi:* nomor pesanan yang tertempel di tas belanja; setiap pelayan
restoran yang menerima tas bisa membacanya tanpa ditanya.
*Di repo:* `order.id` di-set di order, dibaca inventory & payment
(`baggage.order.id`).

**Resource**
"Kartu identitas" aplikasi yang melekat di semua span-nya.
*Analogi:* seragam kerja dengan name-tag — dari jauh tahu ini karyawan
bagian mana. *Di repo:* `service.name=order-service`,
`deployment.environment.name=lab`.

**Attributes**
Keterangan tambahan yang ditempel di satu span.
*Analogi:* detail di satu bab: "produk A123, sisa 10".
*Di repo:* `product.stock`, `payment.amount`, `order.quantity`.

**Events**
Catatan momen penting *di dalam* satu span, lengkap dengan waktunya.
*Analogi:* catatan kaki di bab: "pukul 10.00 mulai memasak, 10.07 selesai".
*Di repo:* `payment_started` → `payment_completed` di span payment.

**Span status (ERROR) & RecordError**
Tanda bahwa langkah itu gagal, plus penyimpanan detail errornya.
*Analogi:* bab ditutup dengan stempel "GAGAL" + lampiran surat keluhan.
*Di repo:* saat chaos 100%, span payment dan parent checkout jadi merah.

**Jaeger**
Aplikasi untuk *melihat* semua rekaman — pencarian + tampilan waterfall.
*Analogi:* ruang arsip CCTV yang bisa diputar per kejadian.
*Di repo:* UI di `localhost:16686`.

**Waterfall**
Tampilan trace seperti tangga: bar horizontal berurutan, makin dalam makin
kecil.
*Analogi:* rincian perjalanan: jalan A (3 mnt) lalu di dalamnya jalan B
(1 mnt). Bar yang panjang = langkah yang lama.

**Context propagation**
Proses "menyalakan-terus" data penting (trace id, baggage) saat request
berpindah fungsi maupun service.
*Analogi:* testi yang harus diserahkan ke petugas berikutnya di tiap
lorong — kalau ada yang lupa, riwayat putus.
*Di repo:* `context.Context` diteruskan di semua handler → client → transport.

---

## Kubernetes & Container

**Container**
Cara mengemas aplikasi + semua kebutuhannya jadi satu paket yang jalan di mana
saja.
*Analogi:* bekal makan siang dalam kotak tertutup — sama rasanya di mana pun
dibuka. *Di repo:* image `otel-shop/order:local`.

**Image**
"Foto lengkap" aplikasi siap-jalan — blueprint dari container.
*Analogi:* cetakan kue vs kuenya: image = cetakan, container = kue yang
dijalankan. *Di repo:* dibangun lewat Dockerfile, multi-stage.

**Multi-stage build**
Cara membangun image dua babak: babak besar untuk men-compile, babak kecil
hanya membawa hasilnya.
*Analogi:* dapur besar untuk memasak, lalu saosnya saja yang dibawa ke meja.
*Di repo:* `golang:1.27-alpine` → `scratch`.

**Image scratch**
Image kosong absolut — tanpa sistem operasi, hanya aplikasinya.
*Analogi:* membawa hanya isinya, tanpa kotaknya. Kecil & minim bahan
kerusakan. *Di repo:* runtime ketiga service.

**Docker vs Podman**
Dua alat untuk menjalankan container — perintahnya mirip, mesinnya beda.
*Analogi:* dua merek motosiklet yang sama-sama bisa ke tujuan.
*Di repo:* build pakai Podman (PRD), tapi k3d hanya mau Docker — Day 1 drama.

**k3d**
Alat untuk membuat Kubernetes mini di laptop, berbasis container.
*Analogi:* stadion mini di meja belajar: semua elemen bola ada, ukurannya
kecil. *Di repo:* cluster `otel-shop` (1 server + 1 agent).

**Cluster**
Kumpulan mesin (nyata/maya) yang dikelola sebagai satu oleh Kubernetes.
*Analogi:* gudang berisi rak; manajer gudang mengatur semuanya.
*Di repo:* `k3d cluster create otel-shop` di `scripts/cluster.sh`.

**Pod**
Satu "kandang" tempat container berjalan — unit terkecil yang dijalankan
Kubernetes.
*Analogi:* satu kandang ayam; bisa berisi satu ayam (container).
*Di repo:* `pod/order-service-xxxx`.

**Deployment**
Pengawas yang memastikan jumlah pod sesuai keinginan dan otomatis mengganti
yang mati.
*Analogi:* mandor yang memastikan selalu ada 1 karyawan di pos; yang sakit
diganti. *Di repo:* `deploy/order-service` → replicas: 1.

**Service**
Alamat tetap untuk menemukan pod (yang alamatnya sering berubah).
*Analogi:* nomor telepon toko — kasirnya bisa berganti, nomornya tetap.
*Di repo:* `order-service`, `inventory-service`, dst.

**NodePort**
Pintu keluar khusus supaya orang dari luar cluster (laptop kita) bisa masuk.
*Analogi:* nomor loker di pintu depan gedung — dari luar tinggal ambil
loker itu. *Di repo:* 18080/18081/18082 + 15432 + 14317/14318 + 16686.

**ClusterIP vs NodePort**
ClusterIP: alamat internal (hanya penghuni gedung). NodePort: punya pintu
dari luar.
*Analogi:* extension telepon kantor vs nomor bisa-dihubungi-dari-rumah.
*Di repo:* jaeger 4317 ClusterIP-ish (diakses collector), 16686 NodePort.

**Namespace**
Kamar terpisah di dalam cluster supaya benda tidak campur aduk.
*Analogi:* rak berlabel per anggota keluarga. *Di repo:* namespace
`otel-shop`.

**ConfigMap**
Tempat menaruh file/aturan yang dibaca pod saat jalan.
*Analogi:* papan pengumuman di dapur: resep ditempel, koki membaca.
*Di repo:* config collector, init.sql Postgres.

**Secret**
ConfigMap versi rahasia — untuk password/token.
*Analogi:* amplop tertutup, bukan papan pengumuman.
*Di repo:* `postgres-credentials` (user/password db).

**Probe (readiness/liveness)**
Alarm kesehatan otomatis: readiness = "siap menerima tamu belum?",
liveness = "masih hidup tidak?".
*Analogi:*resepsionis mengetuk: siap menerima?; satpam mengetuk: masih
sadar? *Di repo:* probe `/health` tiap 5s/10s di 3 service.

**Rollout restart**
Minta Kubernetes mengganti pod lama dengan pod baru — tanpa mengubah
pengaturan.
*Analogi:* ganti shift karyawan: pekerjaan tetap jalan, orangnya berganti.
*Di repo:* `kubectl rollout restart deployment/...` (Day 6: wajib setelah
image diperbarui).

**Scale**
Menambah/mengurangi jumlah salinan pod.
*Analogi:* menambah kursi di restoran saat ramai, mengurangi saat sepi.
*Di repo:* `kubectl scale deploy/otel-collector --replicas=0` (eksperimen
collector down — checkout tetap jalan).

**kubectl**
Alat kendali Kubernetes dari terminal — seperti remote TV.
*Analogi:* satu remote untuk semua: lihat pod, ganti env, restart, scale.
*Di repo:* dipakai di semua script deploy/chaos-test.

---

## Go & Testing

**Go module**
Satu paket kode Go dengan identitas sendiri dan daftar ketergantungannya.
*Analogi:* satu buku dengan nomor ISBN + daftar pustaka.
*Di repo:* 4 module: order, inventory, payment, telemetry.

**go.mod / go.sum**
go.mod = daftar isi + ketergantungan buku; go.sum = checksum agar tidak
ada yang mengelabui.
*Analogi:* daftar bahan + tanda tangan pembuktian bahan asli.

**go.work (workspace)**
File yang memberi tahu Go: "kelola beberapa module sekaligus sebagai satu
proyek".
*Analogi:* rak yang menampung beberapa buku supaya dibaca serentak.
*Di repo:* `go.work` memuat 4 module.

**Replace directive**
Perintah di go.mod: "untuk library X, pakai salinan lokal ini, jangan dari
internet."
*Analogi:* "resep ini pakai bumbu buatan sendiri, bukan yang dibeli."
*Di repo:* `replace github.com/otel-shop/telemetry => ../../pkg/telemetry`
(wajib agar Docker build menemukan module telemetry).

**Interface**
Daftar janji perilaku — tanpa peduli siapa yang melakukannya.
*Analogi:* "petugas kasir harus bisa: scan, terima uang, kembalikan uang".
Tidak peduli manusia atau mesin. *Di repo:* `StockStore`,
`InventoryClient`, `PaymentClient`.

**Mock / fake**
Peniru yang menggantikan benda asli saat latihan (testing).
*Analogi:* boneka layanan sebagai pelanggan untuk latihan kasir.
*Di repo:* `fakeInventory`, `fakeStore` di test.

**httptest**
Alat Go untuk membuat server palsu super ringan hanya untuk test.
*Analogi:* toko mini replika di gudang untuk latihan kasir — tidak
mengganggu toko sungguhan. *Di repo:* `client_test.go`,
`checkout_integration_test.go`.

**sqlmock**
Server database palsu yang bisa diatur jawabannya (baris ada / tidak ada /
error).
*Analogi:* telepon palsu untuk latihan CS — semua skenario jawaban bisa
disetel. *Di repo:* `product_test.go` — test DB tanpa Postgres nyata.

**pgx**
Driver (pengemudi) untuk Go berbicara dengan PostgreSQL.
*Analogi:* penerjemah khusus yang memahami bahasa Postgres.
*Di repo:* `jackc/pgx/v5/stdlib` + `otelsql` membungkusnya.

**context.Context**
"Batas waktu & data bawaan" yang diteruskan ke semua pemanggilan berikutnya.
*Analogi:* kartu barang yang wajib diserahkan ke tiap petugas berikutnya —
kalau ada yang lupa meneruskan, riwayat & batas waktu putus.
*Di repo:* sumber utama bug halus di Go — di sini diteruskan konsisten.

**Unit test vs Integration test vs E2E test**
Unit = uji satu komponen kecil dengan peniru (cepat, terisolasi).
Integration = uji sambungan beberapa komponen asli.
E2E = uji seluruh sistem dari pintu depan (host → cluster).
*Analogi:* uji mesin kamera / uji kamera terhubung perekam / uji seluruh
sistem CCTV dari ruang monitoring. *Di repo:* ketiganya ada
(`checkout_test.go`, `checkout_integration_test.go`, `test.sh`).

**Coverage & coverage gate**
Coverage = persen kode yang pernah dijalankan saat test. Gate = aturan
"kalau kurang dari 70%, tolak".
*Analogi:* persen ruangan yang pernah dikunjungi inspektor; gate = pintu
yang macet kalau persennya kurang. *Di repo:* `check.sh` — hasil
82.1–96.7%.

---

## Tools & Quality

**gofmt**
Perapi otomatis format kode Go — satu gaya untuk semua.
*Analogi:* rapi otomatis: buku selalu tersusun sama.

**go vet**
Dokter ringan yang memeriksa kode tanpa menjalankannya — cari kesalahan
ciri khas.
*Analogi:* skrining kesehatan: cepat, bisa menemukan gejala awal.

**golangci-lint**
Pemeriksa kode yang lebih ketat: gaya, pola berbahaya, kebiasaan buruk.
*Analogi:* audit lengkap vs skrining. *Di repo:* `.golangci.yml` + di
`check.sh`.

**gitleaks**
Anjing pelacak "rahasia" (password, token) di dalam kode.
*Analogi:* detektif yang menyendus amplop tertinggal di kopor sebelum
kopor dikirim. *Di repo:* `gitleaks dir .` sebelum tiap commit penting.

**Chaos engineering**
Seni merusak sistem sengaja (terkontrol) untuk membuktikan sistem kuat dan
untuk *melihat* bagaimana gagal.
*Analogi:* matikan satu generator gedung secara sengaja untuk memastikan
lampu cadangan menyala. *Di repo:* delay/error payment bisa disetel dari
env tanpa mengubah kode.

**CI runner (Continuous Integration)**
Komputer milik layanan (mis. GitHub) yang menjalankan pemeriksaan otomatis
setiap kode di-push.
*Analogi:* penjaga gerbang yang mengetes setiap kiriman sebelum masuk gudang.
*Di repo:* belum dipasang (di luar scope PRD) — unit test sudah dirancang
CI-ready (tanpa DB/cluster).

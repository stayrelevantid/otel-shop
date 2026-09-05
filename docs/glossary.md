# Glossary / Glosarium — OTel-Shop Lab

**🇬🇧 EN** — Term dictionary for beginners. Each term explained *like you're five*, with a daily-life analogy and where it lives in this repo. No reading order needed — use it as a dictionary.

**🇮🇩 ID** — Kamus istilah untuk pemula. Setiap istilah dijelaskan *seperti menjelaskan ke anak kecil*, lengkap dengan analogi sehari-hari dan di mana istilah itu muncul di repo ini. Tidak perlu berurutan — pakai sebagai kamus.

Legend / Keterangan: 🇬🇧 English · 🇮🇩 Bahasa Indonesia · 📍 where it lives in the repo / lokasinya di repo.

---

## 1. Tracing & OpenTelemetry

**OpenTelemetry (OTel)**
🇬🇧 A standard set of tools and rules so every application can "tell what it did" the same way. *Analogy:* a common official language — everyone speaks it, everyone understands.
🇮🇩 Kumpulan alat dan aturan standar supaya semua aplikasi bisa "menceritakan apa yang mereka lakukan" dengan cara yang sama. *Analogi:* bahasa baku — semua pihak pakai, semua paham.
📍 `pkg/telemetry`

**Distributed tracing**
🇬🇧 A way to see one request's journey as it hops across many applications. *Analogy:* a parcel from Jakarta to Solo passing several post offices — tracing records every stop in one history.
🇮🇩 Cara melihat perjalanan satu permintaan yang melompat-lompat antar banyak aplikasi. *Analogi:* paket Jakarta→Solo melewati beberapa kantor pos — tracing mencatat semua perhentian dalam satu riwayat.
📍 one checkout = order → inventory → payment

**Trace**
🇬🇧 The complete history of one request from start to finish. *Analogy:* one book telling one journey from cover to cover.
🇮🇩 Riwayat lengkap satu permintaan dari awal sampai selesai. *Analogi:* satu buku yang menceritakan satu perjalanan dari depan ke belakang.
📍 has one `traceID`, e.g. / contoh `d87f3b8f...`

**Span**
🇬🇧 A note about *one step* inside the journey. *Analogy:* one chapter in the journey book — "left home (3 min)".
🇮🇩 Catatan tentang *satu langkah* di dalam perjalanan. *Analogi:* satu bab dalam buku perjalanan — "berangkat dari rumah (3 menit)".
📍 `POST /checkout`, `check-inventory`, `sql.conn.query` are all spans / semuanya span

**Trace ID**
🇬🇧 A unique number attached to all spans of one request. *Analogy:* a tracking (resi) number — every office holding your parcel writes the same number.
🇮🇩 Nomor unik yang menempel di semua span dari satu permintaan. *Analogi:* nomor resi — semua kantor pos yang memegang paketmu mencatat nomor yang sama.
📍 what joins order+inventory+payment into one trace / yang menyatukan span jadi satu trace

**Span parent–child**
🇬🇧 A span can have a parent — small steps nest inside big ones. *Analogy:* chapter 3 has sub-chapters 3.1, 3.2.
🇮🇩 Span bisa punya orang tua — langkah kecil ada di dalam langkah besar. *Analogi:* bab 3 punya sub-bab 3.1, 3.2.
📍 `sql.conn.query` under / di bawah `GET /inventory/{id}`, under / di bawah `POST /checkout`

**Instrumentation**
🇬🇧 The process of installing "cameras" inside code so activity gets recorded. *Analogy:* mounting CCTV in every room.
🇮🇩 Proses memasang "kamera" di dalam kode supaya aktivitas tercatat. *Analogi:* memasang CCTV di tiap ruangan rumah.
📍 otelhttp & otelsql = auto; span `validate-order` = manual

**Auto vs manual instrumentation**
🇬🇧 Auto = install once, it records everything (middleware). Manual = you choose the important points to record. *Analogy:* automatic door camera vs handwritten diary.
🇮🇩 Otomatis = sekali pasang langsung merekam semua (middleware). Manual = kita pilih titik penting sendiri. *Analogi:* CCTV otomatis di pintu vs catatan tangan di buku harian.
📍 otelhttp/otelsql = auto; `validate-order` dst = manual

**SDK**
🇬🇧 Ready-made software to install those cameras. *Analogy:* a full CCTV kit — camera + recorder + cables.
🇮🇩 Perangkat lunak siap pakai untuk memasang kamera-kamera itu. *Analogi:* paket CCTV lengkap berisi kamera + perekam + kabel.
📍 `go.opentelemetry.io/otel/sdk`

**TracerProvider**
🇬🇧 The manager of all cameras: controls recording and shipping of results. *Analogy:* the building's CCTV monitoring center.
🇮🇩 Manajer semua kamera: mengatur perekaman dan pengiriman hasilnya. *Analogi:* pusat monitoring CCTV gedung.
📍 created / dibuat di `pkg/telemetry.Init()`, stopped via / dimatikan lewat `shutdown()`

**Exporter**
🇬🇧 The worker that carries recordings from the app to storage. *Analogy:* the courier delivering CCTV footage to the archive room.
🇮🇩 Tukang yang mengantar rekaman dari aplikasi ke tempat penyimpanan. *Analogi:* kurir yang mengantar video CCTV ke ruang arsip.
📍 OTLP gRPC exporter → `otel-collector:14317`

**OTLP**
🇬🇧 The standard "shipping language/locket" for those recordings. *Analogy:* a standard parcel wrap every post office understands.
🇮🇩 Bahasa/bungkusan standar untuk mengantar rekaman. *Analogi:* format bungkusan standar yang dipahami semua kantor pos.
📍 `14317` = gRPC door / pintu gRPC, `14318` = HTTP door / pintu HTTP di collector

**Collector**
🇬🇧 The middleman receiving recordings from many apps, then forwarding them. *Analogy:* the central post office: couriers drop off, sorting, onward delivery.
🇮🇩 Perantara yang menerima rekaman dari banyak aplikasi lalu meneruskannya. *Analogi:* kantor pos pusat: semua kurir drop-off, diurutkan, dikirim ke gudang akhir.
📍 deployment `otel-collector`, config / konfigurasi `deploy/otel-collector/configmap.yaml`

**Pipeline (receiver → processor → exporter)**
🇬🇧 The workflow inside the collector: intake door → processing → out door. *Analogy:* a factory: materials in, assembled, product out.
🇮🇩 Jalur kerja di dalam collector: pintu masuk → pengolahan → pintu keluar. *Analogi:* pabrik: bahan masuk, dirakit, produk keluar.
📍 `processors: [batch]` in / di configmap collector

**Batch processor**
🇬🇧 Stores recordings briefly, then ships many at once. *Analogy:* school pick-up — collect the kids first, drive together.
🇮🇩 Menyimpan rekaman sebentar lalu mengirim banyak sekaligus. *Analogi:* antar jemput anak sekolah — kumpulkan dulu, baru berangkat bareng.
📍 `processors: [batch]`

**Sampling**
🇬🇧 Deliberately keeping only some recordings (e.g. 1 in 10) to avoid overload. *Analogy:* CCTV with limited storage — records 1 of every 10 passers-by.
🇮🇩 Sengaja hanya menyimpan sebagian rekaman (misal 1 dari 10) supaya tidak berlebihan. *Analogi:* CCTV storage terbatas — rekam 1 dari 10 lewatan.
📍 `OTEL_TRACES_SAMPLER_ARG=0.1` → ~10% checkouts get a trace / hanya ~10% checkout menghasilkan trace

**ParentBased (sampler)**
🇬🇧 Rule: if a request arrives with a "record/don't record" decision, follow it; a new decision is only made for the very first request. *Analogy:* children follow the parents' decision — one family, one choice.
🇮🇩 Aturan: kalau permintaan datang dengan keputusan "direkam/tidak", ikuti; keputusan baru hanya untuk permintaan paling awal. *Analogi:* anak ikut keputusan orang tua — satu keluarga satu pilihan.
📍 `sdktrace.ParentBased(TraceIDRatioBased(...))` in / di `pkg/telemetry`

**Propagator (W3C TraceContext)**
🇬🇧 The rules for "sticking the tracking number" into HTTP headers as the request moves services. *Analogy:* the courier's handover card, copied at every office.
🇮🇩 Aturan cara "menempelkan nomor resi" ke header HTTP saat request pindah service. *Analogi:* kartu barang bawaan kurir yang wajib di-copy tiap kantor pos.
📍 header `traceparent` automatic / otomatis by otelhttp transport

**Baggage**
🇬🇧 A small bag of tiny data carried by the request across services — every service can read it. *Analogy:* the order number stuck on the shopping bag; every waiter receiving the bag can read it without asking.
🇮🇩 Tas kecil berisi data kecil yang dibawa-bawa request saat pindah service — isinya bisa dibaca semua service. *Analogi:* nomor pesanan yang tertempel di tas; setiap pelayan bisa membacanya tanpa ditanya.
📍 `order.id` set in / di-set di order, read in / dibaca inventory & payment (`baggage.order.id`)

**Resource**
🇬🇧 The app's "identity card" attached to all its spans. *Analogy:* a work uniform with a name tag — you know the department from afar.
🇮🇩 "Kartu identitas" aplikasi yang melekat di semua span-nya. *Analogi:* seragam kerja dengan name-tag — dari jauh tahu ini bagian apa.
📍 `service.name=order-service`, `deployment.environment.name=lab`

**Attributes**
🇬🇧 Extra details attached to one span. *Analogy:* the details in one chapter: "product A123, 10 left".
🇮🇩 Keterangan tambahan yang ditempel di satu span. *Analogi:* detail di satu bab: "produk A123, sisa 10".
📍 `product.stock`, `payment.amount`, `order.quantity`

**Events**
🇬🇧 Notes of important moments *inside* one span, each with its own time. *Analogy:* margin notes in a chapter: "10:00 cooking starts, 10:07 done".
🇮🇩 Catatan momen penting *di dalam* satu span, lengkap dengan waktunya. *Analogi:* catatan kaki di bab: "pukul 10.00 mulai memasak, 10.07 selesai".
📍 `payment_started` → `payment_completed` in / di span payment

**Span status (ERROR) & RecordError**
🇬🇧 A marker that the step failed, plus storing the error details. *Analogy:* the chapter closed with a "FAILED" stamp + an attached complaint letter.
🇮🇩 Tanda bahwa langkah itu gagal, plus penyimpanan detail errornya. *Analogi:* bab ditutup stempel "GAGAL" + lampiran surat keluhan.
📍 chaos 100%: payment span and / dan parent checkout turn red / jadi merah

**Jaeger**
🇬🇧 The application for *viewing* all recordings — search + waterfall display. *Analogy:* the CCTV archive room, playable per incident.
🇮🇩 Aplikasi untuk *melihat* semua rekaman — pencarian + tampilan waterfall. *Analogi:* ruang arsip CCTV yang bisa diputar per kejadian.
📍 UI at / di `localhost:16686`

**Waterfall**
🇬🇧 The trace display like a staircase: horizontal bars in sequence, deeper = smaller. *Analogy:* trip details: road A (3 min), and inside it road B (1 min). A long bar = a long step.
🇮🇩 Tampilan trace seperti tangga: bar horizontal berurutan, makin dalam makin kecil. *Analogi:* rincian perjalanan: jalan A (3 mnt) lalu di dalamnya jalan B (1 mnt). Bar panjang = langkah lama.
📍 Jaeger UI

**Context propagation**
🇬🇧 The act of "keeping data alive" (trace id, baggage) as the request crosses functions and services. *Analogy:* a handover note that must be passed to the next officer in every corridor — one forgetful person breaks the history.
🇮🇩 Proses "menyalakan-terus" data penting (trace id, baggage) saat request berpindah fungsi maupun service. *Analogi:* testi yang wajib diserahkan ke petugas berikutnya di tiap lorong — kalau ada yang lupa, riwayat putus.
📍 `context.Context` passed through / diteruskan di semua handler → client → transport

---

## 2. Kubernetes & Container

**Container**
🇬🇧 A way to pack an app plus everything it needs into one portable box. *Analogy:* a packed lunchbox — tastes the same wherever opened.
🇮🇩 Cara mengemas aplikasi + semua kebutuhannya jadi satu paket yang jalan di mana saja. *Analogi:* bekal makan siang dalam kotak — sama rasanya di mana pun dibuka.
📍 image `otel-shop/order:local`

**Image**
🇬🇧 The complete "ready-to-run snapshot" of an app — the blueprint of containers. *Analogy:* cake mold vs cake: image = mold, container = the running cake.
🇮🇩 "Foto lengkap" aplikasi siap-jalan — blueprint dari container. *Analogi:* cetakan kue vs kuenya: image = cetakan, container = kue yang dijalankan.
📍 built via / dibangun lewat Dockerfile, multi-stage

**Multi-stage build**
🇬🇧 Building an image in two acts: a big stage to compile, a small stage carrying only the result. *Analogy:* a big kitchen to cook, then only the sauce goes to the table.
🇮🇩 Membangun image dua babak: babak besar untuk compile, babak kecil hanya membawa hasilnya. *Analogi:* dapur besar untuk memasak, lalu saosnya saja dibawa ke meja.
📍 `golang:1.27-alpine` → `scratch`

**Scratch image**
🇬🇧 An absolutely empty image — no OS, only the app. *Analogy:* carrying only the contents, no box at all. Tiny and hard to break.
🇮🇩 Image kosong absolut — tanpa sistem operasi, hanya aplikasinya. *Analogi:* membawa hanya isinya, tanpa kotak. Kecil & minim bahan kerusakan.
📍 runtime of / dari ketiga service

**Docker vs Podman**
🇬🇧 Two tools to run containers — similar commands, different engines. *Analogy:* two motorcycle brands — both get you there.
🇮🇩 Dua alat untuk menjalankan container — perintahnya mirip, mesinnya beda. *Analogi:* dua merek motor yang sama-sama bisa ke tujuan.
📍 build uses / build pakai Podman (PRD), but k3d only wants Docker / tapi k3d hanya mau Docker — Day 1 drama

**k3d**
🇬🇧 A tool to create a mini Kubernetes on your laptop, container-based. *Analogy:* a mini stadium on your desk — every element of the game, smaller.
🇮🇩 Alat untuk membuat Kubernetes mini di laptop, berbasis container. *Analogi:* stadion mini di meja belajar: semua elemen bola ada, ukurannya kecil.
📍 cluster `otel-shop` (1 server + 1 agent) in / di `scripts/cluster.sh`

**Cluster**
🇬🇧 A set of machines (real/virtual) managed as one by Kubernetes. *Analogy:* a warehouse of racks; the manager controls all of them.
🇮🇩 Kumpulan mesin (nyata/maya) yang dikelola sebagai satu oleh Kubernetes. *Analogi:* gudang berisi rak; manajer gudang mengatur semuanya.
📍 `k3d cluster create otel-shop` in / di `scripts/cluster.sh`

**Pod**
🇬🇧 One "pen" where a container runs — the smallest unit Kubernetes runs. *Analogy:* one chicken coop; it can hold one chicken (container).
🇮🇩 Satu "kandang" tempat container berjalan — unit terkecil yang dijalankan Kubernetes. *Analogi:* satu kandang ayam; bisa berisi satu ayam (container).
📍 `pod/order-service-xxxx`

**Deployment**
🇬🇧 The supervisor ensuring the pod count matches your wish, auto-replacing the dead. *Analogy:* a foreman ensuring one worker is always at the post; sick ones replaced.
🇮🇩 Pengawas yang memastikan jumlah pod sesuai keinginan dan otomatis mengganti yang mati. *Analogi:* mandor memastikan selalu ada 1 karyawan di pos; yang sakit diganti.
📍 `deploy/order-service` → replicas: 1

**Service**
🇬🇧 A fixed address to find pods (whose addresses change often). *Analogy:* the store's phone number — cashiers rotate, the number stays.
🇮🇩 Alamat tetap untuk menemukan pod (yang alamatnya sering berubah). *Analogi:* nomor telepon toko — kasirnya berganti, nomornya tetap.
📍 `order-service`, `inventory-service`, dst.

**NodePort**
🇬🇧 A special exit door so people outside the cluster (your laptop) can get in. *Analogy:* the locker number at the building's front door — grab that locker from outside.
🇮🇩 Pintu keluar khusus supaya orang dari luar cluster (laptop) bisa masuk. *Analogi:* nomor loker di pintu depan gedung — dari luar tinggal ambil loker itu.
📍 18080/18081/18082 + 15432 + 14317/14318 + 16686

**ClusterIP vs NodePort**
🇬🇧 ClusterIP: internal address (building residents only). NodePort: has an outside door. *Analogy:* office phone extension vs a number reachable from home.
🇮🇩 ClusterIP: alamat internal (hanya penghuni gedung). NodePort: punya pintu dari luar. *Analogi:* extension telepon kantor vs nomor bisa-dihubungi-dari-rumah.
📍 jaeger 4317 internal (collector access / diakses collector), 16686 NodePort

**Namespace**
🇬🇧 A separate room inside the cluster so things don't mix. *Analogy:* labeled shelves per family member.
🇮🇩 Kamar terpisah di dalam cluster supaya benda tidak campur aduk. *Analogi:* rak berlabel per anggota keluarga.
📍 namespace `otel-shop`

**ConfigMap**
🇬🇧 A place to put files/settings that pods read while running. *Analogy:* the kitchen notice board: recipes pinned, cooks read them.
🇮🇩 Tempat menaruh file/aturan yang dibaca pod saat jalan. *Analogi:* papan pengumuman di dapur: resep ditempel, koki membaca.
📍 collector config, / init.sql Postgres

**Secret**
🇬🇧 ConfigMap's secretive sibling — for passwords/tokens. *Analogy:* a sealed envelope, not a notice board.
🇮🇩 ConfigMap versi rahasia — untuk password/token. *Analogi:* amplop tertutup, bukan papan pengumuman.
📍 `postgres-credentials` (user/password db)

**Probe (readiness/liveness)**
🇬🇧 Automatic health alarms: readiness = "ready for guests yet?", liveness = "still alive?". *Analogy:* receptionist knocks: ready to receive?; guard knocks: still conscious?
🇮🇩 Alarm kesehatan otomatis: readiness = "siap menerima tamu belum?", liveness = "masih hidup tidak?". *Analogi:* resepsionis mengetuk: siap menerima?; satpam mengetuk: masih sadar?
📍 probe `/health` every / tiap 5s/10s in / di 3 service

**Rollout restart**
🇬🇧 Ask Kubernetes to swap old pods for new ones — without changing settings. *Analogy:* change work shifts: work continues, the person rotates.
🇮🇩 Minta Kubernetes mengganti pod lama dengan pod baru — tanpa mengubah pengaturan. *Analogi:* ganti shift karyawan: pekerjaan tetap jalan, orangnya berganti.
📍 `kubectl rollout restart deployment/...` (Day 6: required after / wajib setelah image updated / image diperbarui)

**Scale**
🇬🇧 Adding/removing copies of a pod. *Analogy:* adding chairs when the restaurant is busy, removing when quiet.
🇮🇩 Menambah/mengurangi jumlah salinan pod. *Analogi:* menambah kursi saat ramai, mengurangi saat sepi.
📍 `kubectl scale deploy/otel-collector --replicas=0` (collector-down experiment / eksperimen collector down — checkout tetap jalan)

**kubectl**
🇬🇧 The terminal controller for Kubernetes — like a TV remote. *Analogy:* one remote for everything: view pods, change env, restart, scale.
🇮🇩 Alat kendali Kubernetes dari terminal — seperti remote TV. *Analogi:* satu remote untuk semua: lihat pod, ganti env, restart, scale.
📍 used in / dipakai di semua script deploy/chaos-test

---

## 3. Go & Testing

**Go module**
🇬🇧 One Go code package with its own identity and dependency list. *Analogy:* one book with an ISBN + bibliography.
🇮🇩 Satu paket kode Go dengan identitas sendiri dan daftar ketergantungannya. *Analogi:* satu buku dengan nomor ISBN + daftar pustaka.
📍 4 modules / module: order, inventory, payment, telemetry

**go.mod / go.sum**
🇬🇧 go.mod = the book's table of contents + dependencies; go.sum = checksums proving nothing was tampered with.
🇮🇩 go.mod = daftar isi + ketergantungan buku; go.sum = checksum agar tidak ada yang mengelabui.
📍 every module / tiap module

**go.work (workspace)**
🇬🇧 A file telling Go: "manage several modules together as one project". *Analogy:* a shelf holding several books so they're read together.
🇮🇩 File yang memberi tahu Go: "kelola beberapa module sekaligus sebagai satu proyek". *Analogi:* rak menampung beberapa buku supaya dibaca serentak.
📍 `go.work` holds / memuat 4 modules / module

**Replace directive**
🇬🇧 An instruction in go.mod: "for library X, use this local copy, not the internet's". *Analogy:* "this recipe uses homemade spice mix, not store-bought".
🇮🇩 Perintah di go.mod: "untuk library X, pakai salinan lokal ini, jangan dari internet." *Analogi:* "resep ini pakai bumbu buatan sendiri, bukan yang dibeli."
📍 `replace github.com/otel-shop/telemetry => ../../pkg/telemetry` (required for / wajib agar Docker build finds / menemukan the telemetry module)

**Interface**
🇬🇧 A list of promised behaviors — regardless of who performs them. *Analogy:* "a cashier must be able to: scan, accept money, give change" — human or machine, doesn't matter.
🇮🇩 Daftar janji perilaku — tanpa peduli siapa yang melakukannya. *Analogi:* "petugas kasir harus bisa: scan, terima uang, kembalikan uang". Manusia atau mesin, tidak peduli.
📍 `StockStore`, `InventoryClient`, `PaymentClient`

**Mock / fake**
🇬🇧 An imitator standing in for the real thing during practice (testing). *Analogy:* a mannequin customer for cashier practice.
🇮🇩 Peniru yang menggantikan benda asli saat latihan (testing). *Analogi:* boneka layanan sebagai pelanggan untuk latihan kasir.
📍 `fakeInventory`, `fakeStore` in / di tests

**httptest**
🇬🇧 A Go tool to spin up super-light fake servers just for tests. *Analogy:* a mini replica store in the warehouse for cashier practice — the real store is untouched.
🇮🇩 Alat Go untuk membuat server palsu super ringan hanya untuk test. *Analogi:* toko mini replika di gudang untuk latihan kasir — tidak mengganggu toko sungguhan.
📍 `client_test.go`, `checkout_integration_test.go`

**sqlmock**
🇬🇧 A fake database server whose answers you can script (rows exist / no rows / error). *Analogy:* a practice phone for CS training — every answer scenario settable.
🇮🇩 Server database palsu yang bisa diatur jawabannya (baris ada / tidak ada / error). *Analogi:* telepon palsu untuk latihan CS — semua skenario jawaban bisa disetel.
📍 `product_test.go` — DB tests without / test DB tanpa real Postgres / Postgres nyata

**pgx**
🇬🇧 The driver that lets Go speak to PostgreSQL. *Analogy:* a dedicated interpreter who understands Postgres.
🇮🇩 Driver (pengemudi) untuk Go berbicara dengan PostgreSQL. *Analogi:* penerjemah khusus yang memahami bahasa Postgres.
📍 `jackc/pgx/v5/stdlib` + `otelsql` wrapping / membungkus it / itu

**context.Context**
🇬🇧 "The deadline & carried data" passed to all following calls. *Analogy:* the handover card that must reach every next officer — one person forgets, history & deadline break.
🇮🇩 "Batas waktu & data bawaan" yang diteruskan ke semua pemanggilan berikutnya. *Analogi:* kartu barang yang wajib diserahkan ke tiap petugas berikutnya — kalau ada yang lupa, riwayat & batas waktu putus.
📍 source of subtle Go bugs / sumber bug halus di Go — passed consistently / di sini diteruskan konsisten

**Unit vs Integration vs E2E test**
🇬🇧 Unit = test one small component with fakes (fast, isolated). Integration = test the connection between real components. E2E = test the whole system from the front door (host → cluster). *Analogy:* test the camera engine / test camera connected to recorder / test the whole CCTV system from the monitoring room.
🇮🇩 Unit = uji satu komponen kecil dengan peniru (cepat, terisolasi). Integration = uji sambungan beberapa komponen asli. E2E = uji seluruh sistem dari pintu depan (host → cluster). *Analogi:* uji mesin kamera / uji kamera terhubung perekam / uji seluruh sistem dari ruang monitoring.
📍 all three exist / ketiganya ada: `checkout_test.go`, `checkout_integration_test.go`, `test.sh`

**Coverage & coverage gate**
🇬🇧 Coverage = percent of code ever executed during tests. Gate = the rule "below 70% → reject". *Analogy:* percent of rooms the inspector visited; the gate jams if the percent is short.
🇮🇩 Coverage = persen kode yang pernah dijalankan saat test. Gate = aturan "kalau kurang dari 70%, tolak". *Analogi:* persen ruangan yang pernah dikunjungi inspektor; gate = pintu macet kalau persennya kurang.
📍 `check.sh` — result / hasil 82.1–96.7%

---

## 4. Tools & Quality

**gofmt**
🇬🇧 The automatic Go code formatter — one style for everyone. *Analogy:* auto-tidy: books always arranged the same way.
🇮🇩 Perapi otomatis format kode Go — satu gaya untuk semua. *Analogi:* rapi otomatis: buku selalu tersusun sama.
📍 gate in / di `check.sh`

**go vet**
🇬🇧 A light doctor examining code without running it — finds telltale mistakes. *Analogy:* a quick health screening — fast, catches early symptoms.
🇮🇩 Dokter ringan yang memeriksa kode tanpa menjalankannya — cari kesalahan ciri khas. *Analogi:* skrining kesehatan: cepat, bisa menemukan gejala awal.
📍 in / di `check.sh`

**golangci-lint**
🇬🇧 The stricter code examiner: style, dangerous patterns, bad habits. *Analogy:* a full audit vs a quick screening.
🇮🇩 Pemeriksa kode yang lebih ketat: gaya, pola berbahaya, kebiasaan buruk. *Analogi:* audit lengkap vs skrining.
📍 `.golangci.yml` + in / di `check.sh`

**gitleaks**
🇬🇧 A dog that sniffs out "secrets" (passwords, tokens) inside code. *Analogy:* a detective sniffing envelopes left in luggage before shipping.
🇮🇩 Anjing pelacak "rahasia" (password, token) di dalam kode. *Analogi:* detektif yang menyendus amplop tertinggal di kopor sebelum kopor dikirim.
📍 `gitleaks dir .` before / sebelum every important commit / tiap commit penting

**Chaos engineering**
🇬🇧 The art of deliberately (and controllably) breaking a system to prove it is strong and to *watch* how it fails. *Analogy:* switching off one building generator on purpose to prove the backup lights come on.
🇮🇩 Seni merusak sistem sengaja (terkontrol) untuk membuktikan sistem kuat dan untuk *melihat* bagaimana gagal. *Analogi:* matikan satu generator gedung sengaja untuk memastikan lampu cadangan menyala.
📍 payment delay/error settable from / bisa disetel dari env without code changes / tanpa mengubah kode

**CI runner (Continuous Integration)**
🇬🇧 A service-owned computer (e.g. GitHub) that runs automatic checks every time code is pushed. *Analogy:* a gatekeeper testing every delivery before it enters the warehouse.
🇮🇩 Komputer milik layanan (mis. GitHub) yang menjalankan pemeriksaan otomatis setiap kode di-push. *Analogi:* penjaga gerbang yang mengetes setiap kiriman sebelum masuk gudang.
📍 not installed yet / belum dipasang (outside / di luar PRD scope) — unit tests are CI-ready / unit test sudah dirancang CI-ready (no DB/cluster / tanpa DB/cluster)
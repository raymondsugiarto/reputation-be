-- Seed for platform internal admin + sample org owner. The platform internal
-- admin lives on the same organization as the sample org owner for now
-- (it's the only seed org). The two users are distinguished by `user_type`:
--   * INTERNAL_ADMIN -> /api/admin/* (cross-org visibility)
--   * ADMIN          -> org-owner admin (only sees their own org)
--   * FUNDER         -> legacy seed row kept for compatibility
--
-- Three flavours of customer are seeded to exercise every part of the
-- approval workflow without having to sign up by hand:
--   * APPROVED    — `customer1` (password `customer`), used to verify
--                    that approved customers can sign in and reach
--                    /dashboard.
--   * PENDING_*   — `pending1` … `pending5` (password `customer`), used
--                    to populate the admin approval queue. Three are
--                    INDIVIDUAL and two are COMPANY so the type filter
--                    on /admin/approvals has non-empty buckets on both
--                    sides.
--   * REJECTED    — `rejected1` (password `customer`), used so the
--                    history page and the rejection-remark UX have a
--                    realistic row to render.
--
-- New customers created via /api/customer/sign-up always start in
-- PENDING_APPROVAL — same as these seeds.

INSERT INTO organization (id, name, code, origin, created_at, updated_at)
 VALUES ('1', '1', '1', 'localhost:9000', NOW(), NOW());


INSERT INTO "user" (id, organization_id, user_type, created_at, updated_at)
 VALUES
    ('1', '1', 'INTERNAL_ADMIN', NOW(), NOW()),
    ('2', '1', 'FUNDER', NOW(), NOW()),
    ('3', '1', 'CUSTOMER', NOW(), NOW()),
    ('4', '1', 'CUSTOMER', NOW(), NOW()),
    ('5', '1', 'CUSTOMER', NOW(), NOW()),
    ('6', '1', 'CUSTOMER', NOW(), NOW()),
    ('7', '1', 'CUSTOMER', NOW(), NOW()),
    ('8', '1', 'CUSTOMER', NOW(), NOW()),
    ('9', '1', 'CUSTOMER', NOW(), NOW())
;
-- internal admin (username 'admin' / pass 'admin')
-- customer1  (password 'customer')  -> APPROVED
-- pending1..3 (INDIVIDUAL)          -> PENDING_APPROVAL
-- pending4..5 (COMPANY)             -> PENDING_APPROVAL
-- rejected1                         -> REJECTED

INSERT INTO user_credential (id, organization_id, user_id, username, password, created_at, updated_at)
 VALUES
    ('1', '1', '1', 'admin',    '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass admin
    ('2', '1', '2', 'owner',    '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass admin
    ('3', '1', '3', 'customer1','$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass customer (APPROVED)
    ('4', '1', '4', 'pending1', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass customer (PENDING)
    ('5', '1', '5', 'pending2', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass customer (PENDING)
    ('6', '1', '6', 'pending3', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass customer (PENDING)
    ('7', '1', '7', 'pending4', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass customer (PENDING)
    ('8', '1', '8', 'pending5', '$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW()), -- pass customer (PENDING)
    ('9', '1', '9', 'rejected1','$2a$07$rXsQxQHRwxwHYNTzHKTl.eilofdCZ9Ci0TTJmLdV6I7rxsYn/O74.', NOW(), NOW())  -- pass customer (REJECTED)
;

-- ---------------------------------------------------------------------------
-- APPROVED customer — usable for the customer login flow.
-- ---------------------------------------------------------------------------

INSERT INTO customer (id, organization_id, user_id, customer_type, nama_lengkap, status,
                      approved_by, approved_at, created_at, updated_at)
 VALUES ('seed-cust-approved', '1', '3', 'INDIVIDUAL', 'Andi Pelanggan',
         'APPROVED', 'system', NOW() - INTERVAL '30 days', NOW() - INTERVAL '31 days', NOW())
;

-- ---------------------------------------------------------------------------
-- PENDING_APPROVAL customers — populates /admin/approvals queue.
--
-- created_at is staggered across the last few days so the "diperiksa X
-- jam lalu" / "X hari lalu" labels on the queue page render
-- meaningfully. TanggalLahir / PendirianTanggal are realistic so the
-- demo doesn't look obviously synthetic.
-- ---------------------------------------------------------------------------

-- Pending INDIVIDUAL #1 — Budi, daftar 5 menit lalu.
INSERT INTO customer (id, organization_id, user_id, customer_type,
                      nama_lengkap, nomor_ktp, nomor_npwp,
                      tanggal_lahir, kota_lahir,
                      status_pernikahan, pendidikan_terakhir, lama_tinggal,
                      alamat_jalan, kecamatan, kota_kabupaten, provinsi, kode_pos,
                      sama_dengan_domisili,
                      status, created_at, updated_at)
 VALUES ('seed-cust-pending-1', '1', '4', 'INDIVIDUAL',
         'Budi Santoso', '3275012309870001', '12.345.678.9-012.000',
         '1992-04-15', 'Bandung',
         'MARRIED', 'BACHELOR', '5 tahun',
         'Jl. Merdeka No. 17', 'Cibeunying Kidul', 'Kota Bandung', 'Jawa Barat', '40121',
         true,
         'PENDING_APPROVAL', NOW() - INTERVAL '5 minutes', NOW() - INTERVAL '5 minutes')
;

-- Pending INDIVIDUAL #2 — Citra, daftar 2 jam lalu.
INSERT INTO customer (id, organization_id, user_id, customer_type,
                      nama_lengkap, nomor_ktp, nomor_npwp,
                      tanggal_lahir, kota_lahir,
                      status_pernikahan, pendidikan_terakhir, lama_tinggal,
                      alamat_jalan, kecamatan, kota_kabupaten, provinsi, kode_pos,
                      sama_dengan_domisili,
                      alamat_ktp_jalan, kecamatan_ktp, kota_kabupaten_ktp,
                      status, created_at, updated_at)
 VALUES ('seed-cust-pending-2', '1', '5', 'INDIVIDUAL',
         'Citra Lestari', '3174014508910002', '23.456.789.0-123.000',
         '1991-08-05', 'Jakarta',
         'SINGLE', 'MASTER', '3 tahun',
         'Jl. Kemang Utara IX No. 5', 'Bangka', 'Jakarta Selatan', 'DKI Jakarta', '12730',
         false,
         'Jl. Sudirman Kav. 21', 'Karet', 'Jakarta Pusat', -- KTP address differs
         'PENDING_APPROVAL', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours')
;

-- Pending INDIVIDUAL #3 — Dimas, daftar 1 hari lalu.
INSERT INTO customer (id, organization_id, user_id, customer_type,
                      nama_lengkap, nomor_ktp, nomor_npwp,
                      tanggal_lahir, kota_lahir,
                      status_pernikahan, pendidikan_terakhir, lama_tinggal,
                      alamat_jalan, kecamatan, kota_kabupaten, provinsi, kode_pos,
                      sama_dengan_domisili,
                      status, created_at, updated_at)
 VALUES ('seed-cust-pending-3', '1', '6', 'INDIVIDUAL',
         'Dimas Pratama', '3578016512890003', '34.567.890.1-234.000',
         '1989-12-25', 'Surabaya',
         'MARRIED', 'SENIOR_HIGH', '10 tahun',
         'Jl. Mayjend Sungkono Kav. 45', 'Pakis', 'Kota Surabaya', 'Jawa Timur', '60256',
         true,
         'PENDING_APPROVAL', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day')
;

-- Pending COMPANY #1 — PT Nusantara Digital, daftar 30 menit lalu.
--
-- The COMPANY-specific columns use the English snake_case names from
-- migration 000004_customer_company_fields (company_name,
-- establishment_date, business_sector, kbli_code, company_tax_id,
-- shareholders) — the Go model keeps the Indonesian names for
-- readability but the DB columns are English.
INSERT INTO customer (id, organization_id, user_id, customer_type,
                      company_name, establishment_date, business_sector, kbli_code,
                      company_tax_id,
                      alamat_jalan, kecamatan, kota_kabupaten, provinsi, kode_pos,
                      shareholders,
                      status, created_at, updated_at)
 VALUES ('seed-cust-pending-4', '1', '7', 'COMPANY',
         'PT Nusantara Digital', '2018-06-12', 'Teknologi Informasi', '62012',
         '01.234.567.8-901.000',
         'Jl. Jend. Sudirman Kav. 28', 'Setiabudi', 'Jakarta Selatan', 'DKI Jakarta', '12950',
         '[{"id":"sh-1","nama":"Andi Wijaya","jenis":"Perorangan","saham":"70","noKtpNpwp":"3275012345670001","peran":"Direktur Utama"},
           {"id":"sh-2","nama":"Sari Putri","jenis":"Perorangan","saham":"30","noKtpNpwp":"3174014508910002","peran":"Komisaris"}]'::jsonb,
         'PENDING_APPROVAL', NOW() - INTERVAL '30 minutes', NOW() - INTERVAL '30 minutes')
;

-- Pending COMPANY #2 — CV Mitra Sejati, daftar 3 jam lalu.




INSERT INTO customer (id, organization_id, user_id, customer_type,
                      company_name, establishment_date, business_sector, kbli_code,
                      company_tax_id,
                      alamat_jalan, kecamatan, kota_kabupaten, provinsi, kode_pos,
                      shareholders,
                      status, created_at, updated_at)
 VALUES ('seed-cust-pending-5', '1', '8', 'COMPANY',
         'CV Mitra Sejati', '2015-02-20', 'Perdagangan Umum', '47111',
         '02.345.678.9-012.000',
         'Jl. Asia Afrika No. 145', 'Cikawao', 'Kota Bandung', 'Jawa Barat', '40261',
         '[{"id":"sh-3","nama":"Hendra Gunawan","jenis":"Perorangan","saham":"100","noKtpNpwp":"3275015504900003","peran":"Pemilik"}]'::jsonb,
         'PENDING_APPROVAL', NOW() - INTERVAL '3 hours', NOW() - INTERVAL '3 hours')
;

-- ---------------------------------------------------------------------------
-- REJECTED customer — populates /admin/approvals/history so the rejection
-- branch has a visible row. rejected_at is set within the last 24 hours
-- so it also bumps the "rejected today" counter on the admin dashboard.
-- ---------------------------------------------------------------------------

INSERT INTO customer (id, organization_id, user_id, customer_type,
                      nama_lengkap, nomor_ktp,
                      status, rejected_by, rejected_at, remark,
                      created_at, updated_at)
 VALUES ('seed-cust-rejected-1', '1', '9', 'INDIVIDUAL',
         'Erika Wulandari', '3275016707900004',
         'REJECTED', 'system', NOW() - INTERVAL '4 hours',
         'Foto KTP tidak terbaca. Mohon daftar ulang dengan dokumen yang lebih jelas.',
         NOW() - INTERVAL '6 hours', NOW() - INTERVAL '4 hours')
;

# Dockerfile Libraries

## Intel® Tiber™ Broadcast Suite

Libraries components list based on Trivy scan results.

### video_production_image:latest (Ubuntu 22.04 LTS, kernel 5.15)

## Used Libraries Scan Results (Trivy)

```text
Total: 53 (UNKNOWN: 0, LOW: 43, MEDIUM: 10, HIGH: 0, CRITICAL: 0)
```

```text
+--------------------------------------------------------------------------------------------------------------------------------------------+
¦     Library      ¦ Vulnerability  ¦ Severity ¦      Installed Version       ¦                            Title                             ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ coreutils        ¦ CVE-2016-2781  ¦ LOW      ¦ 8.32-4.1ubuntu1.2            ¦ coreutils: Non-privileged session can escape to the parent   ¦
¦                  ¦                ¦          ¦                              ¦ session in chroot                                            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2016-2781                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ gcc-12-base      ¦ CVE-2022-27943 ¦          ¦ 12.3.0-1ubuntu1~22.04        ¦ binutils: libiberty/rust-demangle.c in GNU GCC 11.2 allows   ¦
¦                  ¦                ¦          ¦                              ¦ stack exhaustion in demangle_const                           ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-27943                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ gpgv             ¦ CVE-2022-3219  ¦          ¦ 2.2.27-3ubuntu2.1            ¦ gnupg: denial of service issue (resource consumption) using  ¦
¦                  ¦                ¦          ¦                              ¦ compressed packets                                           ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-3219                    ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libapparmor1     ¦ CVE-2016-1585  ¦ MEDIUM   ¦ 3.0.4-2ubuntu2.3             ¦ In all versions of AppArmor mount rules are accidentally     ¦
¦                  ¦                ¦          ¦                              ¦ widened when ...                                             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2016-1585                    ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libatomic1       ¦ CVE-2022-27943 ¦ LOW      ¦ 12.3.0-1ubuntu1~22.04        ¦ binutils: libiberty/rust-demangle.c in GNU GCC 11.2 allows   ¦
¦                  ¦                ¦          ¦                              ¦ stack exhaustion in demangle_const                           ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-27943                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libc-bin         ¦ CVE-2016-20013 ¦          ¦ 2.35-0ubuntu3.8              ¦ sha256crypt and sha512crypt through 0.6 allow attackers to   ¦
¦                  ¦                ¦          ¦                              ¦ cause a denial of...                                         ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2016-20013                   ¦
+------------------¦                ¦          ¦                              +                                                              ¦
¦ libc6            ¦                ¦          ¦                              ¦                                                              ¦
¦                  ¦                ¦          ¦                              ¦                                                              ¦
¦                  ¦                ¦          ¦                              ¦                                                              ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libdbus-1-3      ¦ CVE-2023-34969 ¦          ¦ 1.12.20-2ubuntu4.1           ¦ dbus: dbus-daemon: assertion failure when a monitor is       ¦
¦                  ¦                ¦          ¦                              ¦ active and a message...                                      ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-34969                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libgcc-s1        ¦ CVE-2022-27943 ¦          ¦ 12.3.0-1ubuntu1~22.04        ¦ binutils: libiberty/rust-demangle.c in GNU GCC 11.2 allows   ¦
¦                  ¦                ¦          ¦                              ¦ stack exhaustion in demangle_const                           ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-27943                   ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libgcrypt20      ¦ CVE-2024-2236  ¦ MEDIUM   ¦ 1.9.4-3ubuntu3               ¦ libgcrypt: vulnerable to Marvin Attack                       ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-2236                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libgssapi-krb5-2 ¦ CVE-2024-26462 ¦          ¦ 1.19.2-2ubuntu0.3            ¦ krb5: Memory leak at /krb5/src/kdc/ndr.c                     ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26462                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26458 ¦ LOW      ¦                              ¦ krb5: Memory leak at /krb5/src/lib/rpc/pmap_rmt.c            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26458                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26461 ¦          ¦                              ¦ krb5: Memory leak at /krb5/src/lib/gssapi/krb5/k5sealv3.c    ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26461                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libharfbuzz0b    ¦ CVE-2023-25193 ¦          ¦ 2.7.4-1ubuntu3.1             ¦ harfbuzz: allows attackers to trigger O(n^2) growth via      ¦
¦                  ¦                ¦          ¦                              ¦ consecutive marks                                            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-25193                   ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libk5crypto3     ¦ CVE-2024-26462 ¦ MEDIUM   ¦ 1.19.2-2ubuntu0.3            ¦ krb5: Memory leak at /krb5/src/kdc/ndr.c                     ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26462                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26458 ¦ LOW      ¦                              ¦ krb5: Memory leak at /krb5/src/lib/rpc/pmap_rmt.c            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26458                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26461 ¦          ¦                              ¦ krb5: Memory leak at /krb5/src/lib/gssapi/krb5/k5sealv3.c    ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26461                   ¦
+------------------+----------------+----------¦                              +--------------------------------------------------------------¦
¦ libkrb5-3        ¦ CVE-2024-26462 ¦ MEDIUM   ¦                              ¦ krb5: Memory leak at /krb5/src/kdc/ndr.c                     ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26462                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26458 ¦ LOW      ¦                              ¦ krb5: Memory leak at /krb5/src/lib/rpc/pmap_rmt.c            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26458                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26461 ¦          ¦                              ¦ krb5: Memory leak at /krb5/src/lib/gssapi/krb5/k5sealv3.c    ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26461                   ¦
+------------------+----------------+----------¦                              +--------------------------------------------------------------¦
¦ libkrb5support0  ¦ CVE-2024-26462 ¦ MEDIUM   ¦                              ¦ krb5: Memory leak at /krb5/src/kdc/ndr.c                     ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26462                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26458 ¦ LOW      ¦                              ¦ krb5: Memory leak at /krb5/src/lib/rpc/pmap_rmt.c            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26458                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-26461 ¦          ¦                              ¦ krb5: Memory leak at /krb5/src/lib/gssapi/krb5/k5sealv3.c    ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-26461                   ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ liblzma5         ¦ CVE-2020-22916 ¦ MEDIUM   ¦ 5.2.5-2ubuntu1               ¦ Denial of service via decompression of crafted file          ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2020-22916                   ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libncurses6      ¦ CVE-2023-45918 ¦ LOW      ¦ 6.3-2ubuntu0.1               ¦ ncurses 6.4-20230610 has a NULL pointer dereference in       ¦
¦                  ¦                ¦          ¦                              ¦ tgetstr in tinf ......                                       ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-45918                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2023-50495 ¦          ¦                              ¦ ncurses: segmentation fault via _nc_wrap_entry()             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-50495                   ¦
+------------------+----------------¦          ¦                              +--------------------------------------------------------------¦
¦ libncursesw6     ¦ CVE-2023-45918 ¦          ¦                              ¦ ncurses 6.4-20230610 has a NULL pointer dereference in       ¦
¦                  ¦                ¦          ¦                              ¦ tgetstr in tinf ......                                       ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-45918                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2023-50495 ¦          ¦                              ¦ ncurses: segmentation fault via _nc_wrap_entry()             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-50495                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libpcre3         ¦ CVE-2017-11164 ¦          ¦ 2:8.39-13ubuntu0.22.04.1     ¦ pcre: OP_KETRMAX feature in the match function in            ¦
¦                  ¦                ¦          ¦                              ¦ pcre_exec.c                                                  ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2017-11164                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libpng16-16      ¦ CVE-2022-3857  ¦          ¦ 1.6.37-3build5               ¦ libpng: Null pointer dereference leads to segmentation fault ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-3857                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libsdl2-2.0-0    ¦ CVE-2022-4743  ¦          ¦ 2.0.20+dfsg-2ubuntu1.22.04.1 ¦ SDL2: memory leak in GLES_CreateTexture() in                 ¦
¦                  ¦                ¦          ¦                              ¦ render/opengles/SDL_render_gles.c                            ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-4743                    ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libsndfile1      ¦ CVE-2022-33064 ¦ MEDIUM   ¦ 1.0.31-2ubuntu0.1            ¦ libsndfile: off-by-one error in function wav_read_header in  ¦
¦                  ¦                ¦          ¦                              ¦ src/wav.c leads to code execution...                         ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-33064                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2021-4156  ¦ LOW      ¦                              ¦ libsndfile: heap out-of-bounds read in src/flac.c in         ¦
¦                  ¦                ¦          ¦                              ¦ flac_buffer_copy                                             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2021-4156                    ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ libssl3          ¦ CVE-2022-40735 ¦ MEDIUM   ¦ 3.0.2-0ubuntu1.15            ¦ The Diffie-Hellman Key Agreement Protocol allows use of long ¦
¦                  ¦                ¦          ¦                              ¦ exponents that arguably...                                   ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-40735                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-2511  ¦ LOW      ¦                              ¦ openssl: Unbounded memory growth with session handling in    ¦
¦                  ¦                ¦          ¦                              ¦ TLSv1.3                                                      ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-2511                    ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-4603  ¦          ¦                              ¦ openssl: Excessive time spent checking DSA keys and          ¦
¦                  ¦                ¦          ¦                              ¦ parameters                                                   ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-4603                    ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-4741  ¦          ¦                              ¦ openssl: Use After Free with SSL_free_buffers                ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-4741                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libstdc++6       ¦ CVE-2022-27943 ¦          ¦ 12.3.0-1ubuntu1~22.04        ¦ binutils: libiberty/rust-demangle.c in GNU GCC 11.2 allows   ¦
¦                  ¦                ¦          ¦                              ¦ stack exhaustion in demangle_const                           ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-27943                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libsystemd0      ¦ CVE-2023-7008  ¦          ¦ 249.11-0ubuntu3.12           ¦ systemd-resolved: Unsigned name response in signed zone is   ¦
¦                  ¦                ¦          ¦                              ¦ not refused when DNSSEC=yes...                               ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-7008                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libtinfo6        ¦ CVE-2023-45918 ¦          ¦ 6.3-2ubuntu0.1               ¦ ncurses 6.4-20230610 has a NULL pointer dereference in       ¦
¦                  ¦                ¦          ¦                              ¦ tgetstr in tinf ......                                       ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-45918                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2023-50495 ¦          ¦                              ¦ ncurses: segmentation fault via _nc_wrap_entry()             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-50495                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libudev1         ¦ CVE-2023-7008  ¦          ¦ 249.11-0ubuntu3.12           ¦ systemd-resolved: Unsigned name response in signed zone is   ¦
¦                  ¦                ¦          ¦                              ¦ not refused when DNSSEC=yes...                               ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-7008                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ libzstd1         ¦ CVE-2022-4899  ¦          ¦ 1.4.8+dfsg-3build1           ¦ zstd: mysql: buffer overrun in util.c                        ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-4899                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ login            ¦ CVE-2023-29383 ¦          ¦ 1:4.8.1-2ubuntu2.2           ¦ shadow: Improper input validation in shadow-utils package    ¦
¦                  ¦                ¦          ¦                              ¦ utility chfn                                                 ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-29383                   ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ ncurses-base     ¦ CVE-2023-45918 ¦          ¦ 6.3-2ubuntu0.1               ¦ ncurses 6.4-20230610 has a NULL pointer dereference in       ¦
¦                  ¦                ¦          ¦                              ¦ tgetstr in tinf ......                                       ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-45918                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2023-50495 ¦          ¦                              ¦ ncurses: segmentation fault via _nc_wrap_entry()             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-50495                   ¦
+------------------+----------------¦          ¦                              +--------------------------------------------------------------¦
¦ ncurses-bin      ¦ CVE-2023-45918 ¦          ¦                              ¦ ncurses 6.4-20230610 has a NULL pointer dereference in       ¦
¦                  ¦                ¦          ¦                              ¦ tgetstr in tinf ......                                       ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-45918                   ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2023-50495 ¦          ¦                              ¦ ncurses: segmentation fault via _nc_wrap_entry()             ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-50495                   ¦
+------------------+----------------+----------+------------------------------+--------------------------------------------------------------¦
¦ openssl          ¦ CVE-2022-40735 ¦ MEDIUM   ¦ 3.0.2-0ubuntu1.15            ¦ The Diffie-Hellman Key Agreement Protocol allows use of long ¦
¦                  ¦                ¦          ¦                              ¦ exponents that arguably...                                   ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2022-40735                   ¦
¦                  +----------------+----------¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-2511  ¦ LOW      ¦                              ¦ openssl: Unbounded memory growth with session handling in    ¦
¦                  ¦                ¦          ¦                              ¦ TLSv1.3                                                      ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-2511                    ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-4603  ¦          ¦                              ¦ openssl: Excessive time spent checking DSA keys and          ¦
¦                  ¦                ¦          ¦                              ¦ parameters                                                   ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-4603                    ¦
¦                  +----------------¦          ¦                              +--------------------------------------------------------------¦
¦                  ¦ CVE-2024-4741  ¦          ¦                              ¦ openssl: Use After Free with SSL_free_buffers                ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2024-4741                    ¦
+------------------+----------------¦          +------------------------------+--------------------------------------------------------------¦
¦ passwd           ¦ CVE-2023-29383 ¦          ¦ 1:4.8.1-2ubuntu2.2           ¦ shadow: Improper input validation in shadow-utils package    ¦
¦                  ¦                ¦          ¦                              ¦ utility chfn                                                 ¦
¦                  ¦                ¦          ¦                              ¦ https://avd.aquasec.com/nvd/cve-2023-29383                   ¦
+--------------------------------------------------------------------------------------------------------------------------------------------+
```

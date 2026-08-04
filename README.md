# TOOLS Management System 

Tools Management system adalah sistem untuk manajemen tools
Tools ini terdisi dari 3 bagian, fe, fe-admin, dan be

## System Requirement

Kebutuhan untuk menjalankan sistem ini:
Hardware:

- CPU Minimal dual Core
- RAM Minimal 1GB
- SSD Min 10GB
- Internet

Sotware:

- Golang Min V 1.26
- Postgresql
- root akses

## Installing System

Instalasi Backend TMS disarankan menggunakan linux server (walapun windows ok sih).

- Misalkan akan di install ke dir /opt. maka buat direktory dulu

``` cd /opt
    sudo mkdir tms
    sudo chown -R $USER:$USER
```

- Clone project
``` git clone https://github.com/cosphi84/tms-be.git ```

- Compile binary dulu
``` cd tms-be
    go mod tidy
    go build -o tms-be cmd/api/main.go
```
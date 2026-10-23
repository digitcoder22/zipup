# Maintainer: digits2 <digits2digits0@gmail.com>
pkgname=zipup-suite
pkgver=1.0.0
pkgrel=1
pkgdesc="A lightweight terminal utility suite providing fast zipup and unpack shortcuts"
arch=('x86_64')
url="https://github.com"
license=('MIT')
provides=('zipup' 'unpack')
conflicts=('zipup' 'unpack')
depends=('glibc')
makedepends=('go')
optdepends=('p7zip: required for enabling .7z archive operations')

source=("git+https://github.com.git#tag=v${pkgver}")
sha256sums=('SKIP')

build() {
    cd "${srcdir}/${pkgname}"
    CGO_ENABLED=0 go build -ldflags="-s -w" -o zipup-suite main.go
}

package() {
    cd "${srcdir}/${pkgname}"
    install -Dm755 "zipup-suite" "${pkgdir}/usr/bin/zipup-suite"
    ln -s zipup-suite "${pkgdir}/usr/bin/zipup"
    ln -s zipup-suite "${pkgdir}/usr/bin/unpack"
}

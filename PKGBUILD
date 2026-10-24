# Maintainer: digits2 <digits2digits0@gmail.com>
pkgname=zipup
pkgver=1.0.0
pkgrel=1
pkgdesc="A lightweight terminal utility suite providing fast zipup and unpack shortcuts"
arch=('x86_64' 'aarch64')
url="https://github.com/digitcoder22/zipup"
license=('MIT')
provides=('unpack')
conflicts=('unpack')
depends=('glibc')
makedepends=('go' 'git')
optdepends=('7zip: required for enabling .7z archive operations')

source=("git+${url}.git#tag=v${pkgver}")
sha256sums=('SKIP')

build() {
    cd "${srcdir}/${pkgname}"
    
    export GOPATH="${srcdir}/gopath"
    export GOCACHE="${srcdir}/gocache"
    
    CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o zipup main.go
}

package() {
    cd "${srcdir}/${pkgname}"
    
    install -Dm755 "zipup" "${pkgdir}/usr/bin/zipup"
    ln -s zipup "${pkgdir}/usr/bin/unpack"
    install -Dm644 LICENSE "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}

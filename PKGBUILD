# Maintainer: Lukas Mairock <lukas.mairock@luks.cat>
pkgname=fog
pkgver=1.0.0
pkgrel=1
pkgdesc="System fetcher in Go"
arch=('x86_64')
license=('GPL')
depends=('glibc')
makedepends=('go' 'git')
source=("$pkgname-$pkgver.tar.gz")
noextract=()
sha256sums=('SKIP')

build() {
  cd "$srcdir"
  go build -o "$pkgname" ./main.go
}

package() {
  cd "$srcdir"
  install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
}

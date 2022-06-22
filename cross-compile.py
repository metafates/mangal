import hashlib
import json
import os
import shutil
import subprocess

VERSION = "1.5.1"
TAG = f"v{VERSION}"
DESCRIPTION = "The ultimate CLI manga downloader"
GITHUB = "https://github.com/metafates/mangal"
LICENSE = "MIT"
BIN = "mangal"
GOMOD = "github.com/metafates/mangal"
PLATFORMS = {
    "windows": {
        "arch": ["arm", "arm64", "amd64", "386"],
        "ext": ".exe"
    },
    "linux": {
        "arch": ["arm", "arm64", "amd64", "386"],
        "ext": ""
    },
    "darwin": {
        "arch": ["arm64", "amd64"],
        "ext": ""
    }
}


def sha256(filename):
    chunks = []
    with open(filename, "rb") as f:
        # Read and update hash string value in blocks of 4K
        for byte_block in iter(lambda: f.read(4096), b""):
            hashlib.sha256().update(byte_block)
        chunks.append(hashlib.sha256().hexdigest())
    return "".join(chunks)


def compile_for(goos):
    os.environ["GOOS"] = goos

    ext = PLATFORMS[goos]["ext"]

    for arch in PLATFORMS[goos]["arch"]:
        os.environ["GOARCH"] = arch
        target = os.path.join("bin", f"{BIN}-{goos}-{arch}{ext}")

        subprocess.call([
            "go", "build",
            "-o", target,
            "-ldflags=-s -w",
            GOMOD
        ])


def generate_scoop_manifest():
    win32bin = f"{BIN}-windows-386.exe"
    win64bin = f"{BIN}-windows-amd64.exe"

    manifest = {
        "version": VERSION,
        "description": DESCRIPTION,
        "homepage": GITHUB,
        "license": LICENSE,
        "architecture": {
            "32bit": {
                "hash": sha256(os.path.join("bin", win32bin)),
                "url": f"{GITHUB}/releases/download/{TAG}/{win32bin}",
                "bin": [[win32bin, BIN]],
            },
            "64bit": {
                "hash": sha256(os.path.join("bin", win64bin)),
                "url": f"{GITHUB}/releases/download/{TAG}/{win64bin}",
                "bin": [[win64bin, BIN]],
            }
        },
        "checkver": {"url": GITHUB,
                     "regex": f"{BIN.title()} ([\\d.]+)",
                     },
        "autoupdate": {
            "32bit": {
                "url": f"{GITHUB}/releases/download/{TAG}/{win32bin}"
            },
            "64bit": {
                "url": f"{GITHUB}/releases/download/{TAG}/{win64bin}"
            }
        }
    }

    with open(os.path.join("bin", f"{BIN}.json"), 'w') as scoopFile:
        json.dump(manifest, scoopFile, indent=4)


def generate_homebrew_formula():
    formula = f"""
class {BIN.title()} < Formula
    desc "{DESCRIPTION}"
    homepage "{GITHUB}"
    url "{GITHUB}", :using => :git, :tag => "{TAG}"
    version "{VERSION}"
    sha256 "92da0f4a880f86a5d782f47cf912f0206e2d49c5fcc27d57931ce0ef96e85029"
    license "{LICENSE}"

    depends_on "go" => :build

    def install
        system "go", "build", *std_go_args(ldflags: "-s -w")
    end

    test do
        system "true"
    end
end
    """

    with open(os.path.join("bin", f"{BIN}.rb"), 'w') as f:
        f.write(formula)


# Was not tested
def generate_pkgbuild():
    go_version = "go>=1.8.0"
    pkbuild = f"""
# Maintainer: metafates <fates@duck.com>
pkgname={BIN}
pkgver={VERSION}
pkgrel=1
pkgdesc="{DESCRIPTION}"
arch=('i686' 'x86_64')
url="{GITHUB}"
license=('{LICENSE}')
depends=('glibc')
makedepends=('{go_version}')
checkdepends=('{go_version}')
source=("{GITHUB}/archive/refs/tags/v{VERSION}.tar.gz")

prepare() {{
    cd "{BIN}-{VERSION}"
    mkdir -p build/
}}

build() {{
    cd "{BIN}-{VERSION}"
    export CGO_CPPFLAGS="${{CPPFLAGS}}"
    export CGO_CFLAGS="${{CFLAGS}}"
    export CGO_CXXFLAGS="${{CXXFLAGS}}"
    export CGO_LDFLAGS="${{LDFLAGS}}"
    export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
    go get -d ./...
    go build -o build ./...
}}

package() {{
    cd "{BIN}-{VERSION}"
    install -Dm755 build/{BIN} "{BIN}"/usr/bin/{BIN}
}}
"""

    with open(os.path.join("bin", "PKGBUILD"), 'w') as f:
        f.write(pkbuild)


def main():
    shutil.rmtree("bin")
    for goos in PLATFORMS:
        print(f"Compiling for {goos.title()}")
        compile_for(goos)

    print("Generating Scoop Manifest")
    generate_scoop_manifest()

    print("Generating Homebrew Formula")
    generate_homebrew_formula()

    print("Generating PKGBUILD")
    generate_pkgbuild()

    print("Done")


if __name__ == "__main__":
    main()

import json
import os
import shutil
import subprocess

VERSION = "1.5.2"
TAG = f"v{VERSION}"
DESCRIPTION = "The ultimate CLI manga downloader"
GITHUB = "https://github.com/metafates/mangal"
LICENSE = "MIT"

# Required go version
GO_VERSION = ">=1.8.0"

# Binary name
BIN = "mangal"

# Url of go module
GOMOD = "github.com/metafates/mangal"

# Platforms to compile for
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
    },
    "android": {
        "arch": ["arm64"],
        "ext": ""
    }
}


def sha256(filename):
    """
    Generate SHA256 hash for file
    :param filename: file name
    :return: SHA256 hash
    """

    # call for shasum
    shasum = subprocess.check_output(["shasum", "-a", "256", filename])
    return shasum.split()[0].decode("utf-8")


def compile_for(goos):
    """
    Compile for a given platform
    :param goos: platform to compile for
    """
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
    """
    Generate Scoop manifest for Scoop
    Note: No support for Windows ARM yet since scoop doesn't support it
    """
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
        # scoop special case
        # This will try to match the tag with \/releases\/tag\/(?:v|V)?([\d.]+)
        "checkver": "github",
        "autoupdate": {
            "32bit": {
                "url": f"{GITHUB}/releases/download/v$version/{win32bin}",
                "hash": {
                    "url": f"{GITHUB}/releases/download/v$version/{win32bin}.sha256",
                },
            },
            "64bit": {
                "url": f"{GITHUB}/releases/download/v$version/{win64bin}",
                "hash": {
                    "url": f"{GITHUB}/releases/download/v$version/{win64bin}.sha256",
                },
            }
        }
    }

    with open(os.path.join("bin", f"{BIN}.json"), 'w') as scoopFile:
        json.dump(manifest, scoopFile, indent=4)


def generate_homebrew_formula():
    """
    Generate Homebrew formula for Homebrew
    """
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


def generate_pkgbuild():
    """
    Generate PKGBUILD file for pkgbuild.org
    """
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
makedepends=('go {GO_VERSION}')
checkdepends=('go {GO_VERSION}')
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


def generate_deb_package_for(architecture):
    """
    Generates a deb package for the given architecture.
    :param architecture: The architecture to generate the deb package for.
    """

    if architecture == "386":
        metadata_architecture = "i386"
    elif architecture == "amd64":
        metadata_architecture = "amd64"
    elif architecture == "arm64":
        metadata_architecture = "arm64"
    elif architecture == "arm":
        metadata_architecture = "armhf"
    else:
        raise Exception(f"Unknown architecture: {architecture}")

    metadata = f"""
Package: {BIN}
Version: {VERSION}
Section: utils
Priority: optional
Architecture: {metadata_architecture}
Description: {DESCRIPTION}
Maintainer: metafates <fates@duck.com>
"""

    major_version = VERSION.split(".")[0]
    minor_version = VERSION.split(".")[1]
    revision = VERSION.split(".")[2]

    package_name = f"{BIN}_{major_version}.{minor_version}-{revision}_{metadata_architecture}"

    # make directory with package name at bin folder
    os.mkdir(os.path.join("bin", package_name))

    # make usr/local/bin folder at bin/package_name folder
    os.makedirs(os.path.join("bin", package_name, "usr", "local", "bin"))

    # copy binary to usr/local/bin folder
    shutil.copy(os.path.join("bin", f"{BIN}-linux-{architecture}"),
                os.path.join("bin", package_name, "usr", "local", "bin", BIN))

    # make DEBIAN folder at bin/package_name folder
    os.mkdir(os.path.join("bin", package_name, "DEBIAN"))

    # move metadata to bin/package_name/DEBIAN/control
    with open(os.path.join("bin", package_name, "DEBIAN", "control"), 'w') as f:
        f.write(metadata)

    # build deb package with architecture in name
    subprocess.call(["dpkg-deb", "-b", os.path.join("bin", package_name)])

    # remove bin/package_name folder
    shutil.rmtree(os.path.join("bin", package_name))


def generate_deb_packages():
    # check if dpkg-deb is installed
    try:
        subprocess.check_output(["dpkg-deb", "--version"])
    except subprocess.CalledProcessError:
        print("dpkg-deb is not installed. Please install dpkg-deb to generate a deb package.")
        print("see https://command-not-found.com/dpkg-deb")
        return

    for arch in PLATFORMS["linux"]["arch"]:
        generate_deb_package_for(arch)


def generate_checksums():
    """
    Generate checksums for all files in bin folder
    """

    # check is shasum is installed
    try:
        subprocess.check_output(["shasum", "--version"])
    except subprocess.CalledProcessError:
        print("shasum is not installed. Please install shasum to generate checksums.")
        print("see https://command-not-found.com/shasum")
        return

    # generate shasum for all files in bin folder
    for file in os.listdir(os.path.join("bin")):
        skip_conditions = [
            file.endswith(".json"),
            file.endswith(".sha256"),
            file == "PKGBUILD",
            file.endswith(".rb")
        ]

        if any(skip_conditions):
            continue

        checksum = sha256(os.path.join("bin", file))

        with open(os.path.join("bin", f"{file}.sha256"), 'w') as f:
            f.write(f"{checksum}  {os.path.join('.', file)}")


def main():
    # check if user is in the root directory

    # check if user is in the same directory as this script
    if os.getcwd() != os.path.dirname(os.path.realpath(__file__)):
        print("Please run this script from the root directory of the project.")
        return

    # delete bin folder if it exists
    if os.path.exists(os.path.join("bin")):
        shutil.rmtree(os.path.join("bin"))

    # Cross compile for all platforms listed in PLATFORMS
    for goos in PLATFORMS:
        print(f"Compiling for {goos.title()}")
        compile_for(goos)

    print("Generating Scoop Manifest")
    generate_scoop_manifest()

    print("Generating Homebrew Formula")
    generate_homebrew_formula()

    print("Generating PKGBUILD")
    generate_pkgbuild()

    print("Generating Debian Packages")
    generate_deb_packages()

    print("Generating Checksums")
    generate_checksums()

    print("Done!")


if __name__ == "__main__":
    main()

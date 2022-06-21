import hashlib
import json
import os
import shutil
import subprocess

VERSION = "1.5.1"
TAG = f"v{VERSION}"
DESCRIPTION = "The ultimate CLI manga downloader"
GITHUB = "https://github.com/metafates/mangal"
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
            GOMOD
        ])


def make_scoop_manifest():
    win32bin = f"{BIN}-windows-386.exe"
    win64bin = f"{BIN}-windows-amd64.exe"

    manifest = {
        "version": VERSION,
        "description": DESCRIPTION,
        "homepage": GITHUB,
        "license": "MIT",
        "architecture": {
            "32bit": {
                "hash": sha256(os.path.join("bin", win32bin)),
                "url": f"{GITHUB}/releases/download/{TAG}/{win32bin}",
                "bin": [[win32bin, "mangal"]],
            },
            "64bit": {
                "hash": sha256(os.path.join("bin", win64bin)),
                "url": f"{GITHUB}/releases/download/{TAG}/{win64bin}",
                "bin": [[win64bin, "mangal"]],
            }
        },
        "checkver": {"url": GITHUB,
                     "regex": f"{BIN.title()} ([\\d.]+)",
                     },
        "autoupdate": {
            "32bit": {
                "url": f"{GITHUB}/releases/download/v$version/{BIN}-windows-386.exe"
            },
            "64bit": {
                "url": f"{GITHUB}/releases/download/v$version/{BIN}-windows-amd64.exe"
            }
        }
    }

    with open(os.path.join("bin", "mangal.json"), 'w') as scoopFile:
        json.dump(manifest, scoopFile, indent=4)


def main():
    shutil.rmtree("bin")
    for goos in PLATFORMS:
        print(f"Compiling for {goos.title()}")
        compile_for(goos)

    print("Generating scoop manifest")
    make_scoop_manifest()

    print("Done")


if __name__ == "__main__":
    main()

import os
import subprocess


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


def main():
    for goos in PLATFORMS:
        print(f"Compiling for {goos.title()}")
        compile_for(goos)

if __name__ == "__main__":
    main()

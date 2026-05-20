#!/usr/bin/env python3
from pathlib import Path
from shutil import copyfile


ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / "原型" / "screen.png"
OUTPUT = ROOT / "wx_login" / "screenshots" / "preview" / "project-start-success.png"


def main():
    OUTPUT.parent.mkdir(parents=True, exist_ok=True)
    copyfile(SOURCE, OUTPUT)
    print(f"Saved {OUTPUT.relative_to(ROOT)}")


if __name__ == "__main__":
    main()

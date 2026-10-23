# zipup-suite

A lightweight, high-performance terminal utility suite written in Go providing fast `zipup` and `unpack` syntax shortcuts.

## Installation (Arch Linux)
Note, account registration is temporarily closed, so for now please do the manual install.
```bash
yay -S zipup-src
```
If using paru:
```bash
paru -S zipup-src
```
## Manual Install
```bash
https://github.com/digitcoder22/zipup-src.git
cd zipup-src
makepkg -si
```
## Usage

### Zipup
Compress files or folders quickly by choosing a format and compression level (1-9):
```bash
zipup [format] [level] [file/folder] <optional_destination>
```
*Supported formats: `zip`, `7z`, `tar.gz`, `tar.xz`*

### Unpack
Extract any supported archive automatically to your current directory or a custom path:
```bash
unpack [archive_file] <optional_destination>
```

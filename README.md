# Introduction

This reporisotory provides an experimental alternative [AppImage](https://github.com/probonopd/AppImageKit) runtime (image loader) that does not depend on system-provided glibc, glib, and libfuse to mount and run the image.  Original runtime is embedded into ISO image; this runtime is the "extractor" of a self-extracting ZIP archive.

# Usage

```sh
go get github.com/orivej/static-appimage/...
make-static-appimage APPDIR DESTINATION
```

`APPDIR` must already contain an `AppRun`.

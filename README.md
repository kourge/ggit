ggit
====

A Git library written in Go. This is still in its early development stage, but
the packages are roughly organized in this manner:

* `core` includes types that allow manipulating primitive Git types, such as
  objects (blobs, trees, commits, tags), refs, symbolic refs, etc.
* `config` contains types that read, write, and change INI-style config
  files used by Git.
* `format` provides types that load, save, and manipulate on-disk metadata
  formats like index files and pack files.
* `plumbing` encompasses low-level functions that can be used to directly
  manipulate various Git objects within a repository. (e.g. `hash-object`,
  `cat-file`)
* `porcelain` comprises high-level actions that are assembled from plumbing.
  (e.g. `add`, `commit`, `reset`)

Documentation is available from
[godoc.org](http://godoc.org/github.com/kourge/ggit).

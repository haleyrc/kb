# kb

This repository contains a flat-file "database" of sorts to hold various bits of information I collect
throughout any given day. Eventually, the goal is to migrate this data to an actual database, but
thanks to the limits of available time, this will have to do for now.

## CLI

There is a CLI provided with this repo to make it easier to search and modify the files on disk. In
order to use the CLI you will need Go 1.16+. You will also need to build the tool from within the `kb`
directory, but run it from the root of the repository for the paths to all line up correctly. For example:

```bash
# From the root
cd kb
go build
cd ..
./kb/kb.exe links search -tag=dev
```

## Sections

[Links](links.yaml)
[Shopping](shopping.yaml)

# go-note
this package is note that can be stored online.

support
- storage
    - firebase storage
- editor
    - vim

## Demo
![](https://user-images.githubusercontent.com/7113786/75147434-9fdff700-5740-11ea-802c-d91b875266dd.gif)
## Installation
```bash
go get github.com/yakipote/gostnote
```

edit config toml(~/gostnote/config.toml)
```bash
[Firebase]
keyPath = "your_firestore_key.json"
storageBucket = "your-fire-store-bucket.com"

```

## Features

create new memo
```bash
gostnote new-memo-title
```

show memo list
```bash
gostnote -l
```

show list and edit
```bash
gostonote
```
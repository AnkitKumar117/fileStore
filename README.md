# fileStore
A simple cli and server application in go to implement fileStore.

## Prerequisite

Please [install and set-up Golang](https://go.dev/doc/install) on your system in advance.

## How to run this project?

1. Clone this Project and Navigate to the folder.

```bash
git clone https://github.com/AnkitKumar117/fileStore
cd fileStore
```

2. Build the project using following command.

```bash
go build ./...
```

3. Run the server.

```bash
./server
```

4. You can now use store executable to perfrom following actions.

```bash
./store ls # To list all the files in store.
./store wc # To get word count of all the files in store.
./store update <fileName> # To upsert file in the store.
./store remove <fileName> # To remove file from the store.
./store add <fileName> <fileName> ... # To add multiple files in the store.
./store freq-words [--limit|-n 10] [--order=dsc|asc] # Least or most frequent words in the files.
```
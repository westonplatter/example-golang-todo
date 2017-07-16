# example-golang-todo

Simple Todo List web app.

**The backend is written in Go using the Standard
Library, only**.

The frontend was copy and pasted from the TodoMVC project (backbone).

Pull requests are welcomed and **encouraged**!

## Getting started

Steps for getting up and running,

1. Install go

    See https://golang.org/doc/install

2. Clone the repo

    ```
    git clone https://github.com/westonplatter/example-golang-todo.git
    ```

3. Install project dependencies

    ```
    cd example-golang-todo
    go get
    ```

4. Setup a database

    The project expects a MySQL sever to be accessible via,

    ```sh
    host      = localhost
    username  = root
    password  = (EMPTY)
    ```

    Create a database called `golang_todo_dev`,

    ```sql
    create database golang_todo_dev;
    ```

    Create a table called `Todo`,

    ```sql
    CREATE TABLE `Todo` (
      `Id`          int(11) NOT NULL,
      `Title`       varchar(255) DEFAULT NULL,
      `Category`    varchar(255) DEFAULT NULL,
      `State`       varchar(255) DEFAULT NULL,
      PRIMARY KEY (`Id`)
    ) ENGINE=MyISAM DEFAULT CHARSET=utf8;
    ```

5. Run the web app

    ```sh
    go run server.go
    ```

    Visit [localhost:3000](localhost:3000)

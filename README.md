# example-golang-todo

## getting started

1. Install go

See https://golang.org/doc/install

2. Clone the repo

    git clone https://github.com/westonplatter/example-golang-todo.git

3. Install project dependencies

    cd example-golang-todo
    go get

4. Setup a database

Make sure you have a database called `golang_todo_dev`, with the following table,

```sql
create database Todo;
```

```sql
CREATE TABLE `Todo` (
  `Id`          int(11) NOT NULL,
  `Title`       varchar(255) DEFAULT NULL,
  `Category`    varchar(255) DEFAULT NULL,
  `State`       varchar(255) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
```

4. Run the web app

    go run server.go

Visit [localhost:3000](localhost:3000)

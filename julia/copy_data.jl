using SQLite
using DataFrames
using LibPQ
using BenchmarkTools


function main()
    sqlite_db = SQLite.DB("fakedata.db")
    pg_conn = LibPQ.Connection("dbname=postgres user=postgres password=postgres host=localhost")
    user_df = DBInterface.execute(sqlite_db, "SELECT * FROM users") |> DataFrame

    execute(pg_conn, "DROP TABLE IF EXISTS users")
    execute(
        pg_conn,
        """
        CREATE TABLE users (
            user_id INTEGER PRIMARY KEY,
            user_name TEXT,
            user_age INTEGER,
            user_address TEXT
        )
        """
    )

    execute(pg_conn, "BEGIN;")

    LibPQ.load!(
        user_df,
        pg_conn,
        """
        INSERT INTO users (user_id, user_name, user_age, user_address)
        VALUES            (    \$1,       \$2,      \$3,          \$4)
        """
    )

    execute(pg_conn, "COMMIT;")

    close(pg_conn)
end


@benchmark main()

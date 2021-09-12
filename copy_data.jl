using SQLite
using LibPQ
using BenchmarkTools


function main()
    sqlite_db = SQLite.DB("fakedata.db")
    pg_conn = LibPQ.Connection("dbname=pgdb user=robb password=foob host=localhost")
    rows = DBInterface.execute(sqlite_db, "SELECT * FROM users")

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

    for row in rows
        user_id, user_name, user_age, user_address = row
        LibPQ.load!(
            (user_id = [row.user_id],
             user_name = [row.user_name],
             user_age = [row.user_age],
             user_address = [row.user_address]),
            pg_conn,
            """
            INSERT INTO users (user_id, user_name, user_age, user_address)
            VALUES            (    \$1,       \$2,      \$3,          \$4)
            """
        )
    end

    execute(pg_conn, "COMMIT;")

    close(pg_conn)
end


@benchmark main()

import sqlite3
import random
import time

from faker import Faker


if __name__ == "__main__":
    fake = Faker()

    with sqlite3.connect("fakedata.db") as con:
        con.execute("DROP TABLE IF EXISTS users")
        con.execute(
            """
            CREATE TABLE users (
                user_id INTEGER PRIMARY KEY,
                user_name TEXT,
                user_age INTEGER,
                user_address TEXT
            )
            """
        )
        start = time.time()
        con.executemany(
            "INSERT INTO users (user_id, user_name, user_age, user_address) "
            "VALUES (?, ?, ?, ?)",
            (
                (i, fake.name(), random.randint(10, 100), fake.address())
                for i in range(1_000_000)
            ),
        )
        end = time.time()
    con.close()
    print(f"Finished in {end - start:.1f} seconds.")

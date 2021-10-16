import time

import pandas as pd
import sqlalchemy as sa


if __name__ == "__main__":
    pg_engine = sa.create_engine(
        "postgresql+psycopg2://postgres:postgres@127.0.0.1/postgres"
    )
    chunks = pd.read_sql(
        "SELECT * FROM users",
        "sqlite:///fakedata.db",
        chunksize=10_000,
        index_col="user_id",
    )

    start = time.time()
    with pg_engine.connect() as pg_conn:  # type: ignore
        pg_conn.execute("DROP TABLE IF EXISTS users")
        for df in chunks:
            df.to_sql("users", pg_conn, if_exists="append")
    end = time.time()

    print(f"Finished in {end - start:.1f} seconds")

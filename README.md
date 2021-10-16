SQLite to PostgreSQL
====================

Here are some attempts at timing a transfer of 1MM rows of fake data from a sqlite
database to a postgres database running on `localhost`. I encourage anyone who knows
how to properly write one of these tests to contribute such a script, given it is
sufficiently small (~150 LOC max).


Compute
-------

All of these tests are run on a moderately powered single node:

```
OS: Ubuntu 20.04.3 (Focal Fossa)
Motherboard: PRIME X299-DELUXE
Disk: 256GB SAMSUNG MZNTE256
CPU: Intel(R) Core(TM) i7-7820X CPU @ 3.60GHz
Memory: 64GiB (4 x 16GiB DIMM DDR4 Synchronous 2400 MHz (0.4 ns))
```


Results
-------

Before running commands it's assumed that postgres is already running locally and
that we can switch to the postgres user:

```
sudo su postgres
```

I have both user and password for the `postgres` user set to `postgres` in the
scripts. That usually requires this line:

```
ALTER USER postgres PASSWORD 'postgres';
```

**COPY FROM**

Using pandas, I created a `csv` version of the database for copying, which took a few
seconds, but once that's done the postgres copy is very quick. The following takes
~3-4 seconds:

```
CREATE TABLE IF NOT EXISTS users (
    user_id INT PRIMARY KEY,
    user_name TEXT,
    user_age INT,
    user_address TEXT
);
COPY user_table FROM '/path/to/fakedata.csv' DELIMITER ',' CSV HEADER;
```

The pandas to create the csv:
```python3
import pandas as pd
df = pd.read_sql("users", "sqlite:///fakedata.db", index_col="user_id")
df.to_csv("fakedata.csv")
```


**Python + pandas**

Pandas makes batch processing very easy, so that script will load chunks of rows at a
time, so it's memory efficient and pretty fast.

```
$ python --version
Python 3.9.5
$ python ./python/copy_using_pandas.py
Finished in 26.1 seconds
```

**Julia + LibPQ**

I think that the julia version would require manually chunking the data, so this
version loads all data into memory before writing. This uses a dataframe buffer much
like the pandas version, but is not as time efficient as I had hoped.

```
$ julia --version
julia version 1.6.2
$ julia
julia> ]
(@v1.6) pkg> activate julia/
(julia) pkg> instantiate
(julia) pkg> <BS>
julia> include("julia/copy_data.jl")
BenchmarkTools.Trial: 1 sample with 1 evaluation.
 Single result which took 54.524 s (3.18% GC) to evaluate,
 with a memory estimate of 5.42 GiB, over 107996251 allocations.
```

**Go**

Taking advantage of PostgreSQL support of concurrent writes, we can partition the
original table into chunks (using LIMIT and OFFSET), and have each goroutine handle
copying just that chunk.

Doing it this way lends to a considerable speedup:

```
$ cd golang
$ go run .
Finished in 6.55 seconds
```


Other Languages I'd Like to Try
--------------------------------

I have little experience in the following, but would like to see how they compare.

1. rust
2. C/C++
3. haskell
4. java


Development Environment
---------------

Required installation steps to recreate the tests.

Install the python virtual environment:

```bash
python3.9 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Create the fake sqlite data:
```bash
./create_data.py
```

The julia steps shown above will perform all the environment management you need for it, you just
need [julia 1.6.2](https://julialang.org/downloads)


For `go`, just follow [the installation instructions](https://golang.org/doc/install).

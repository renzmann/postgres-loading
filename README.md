SQLite to PostgreSQL Data Movement
==================================

Here are some attempts at timing a transfer of 1MM rows of fake data from a sqlite
database to a postgres database running on `localhost`. It is important to note that
the runtimes are inversely proportional to my level of proficiency in that
language/framework. I encourage anyone who knows how to properly write one of these
tests to contribute such a script, given it is sufficiently small (~150 LOC max).


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

My very naive (and probably very bad) go version is limited to only 10,000 rows, but
still took 28.4 seconds. A version limited to 1,000 rows ran in about 2.9 seconds,
and a 5,000 row version in around 14 seconds so assuming a linear relationship here 
we can expect 1 MM rows to take about **46 minutes**.

```
$ cd golang
$ go run .
```

Other Languages I'd Like to Try
--------------------------------

I have little experience in the following, but would like to see how they compare.

1. rust
2. C/C++
3. haskell


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

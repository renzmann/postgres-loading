SQLite to PostgreSQL Data Movement
==================================

Here are some attempts at timing a transfer of 1MM rows of fake data from a sqlite
database to a postgres database running on `localhost`. It is important to note that
I am a far less proficient of a programmer in frameworks other than python/pandas. I
encourage anyone who knows how to properly write one of these tests to contribute
such a script, given it is sufficiently small.

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

```bash
python ./python/copy_using_pandas.py
Finished in 26.1 seconds
```

**Julia + LibPQ**

I think that the julia version would require manually chunking the data, so this
version loads all data into memory before writing.

```bash
julia
julia> ]
(@v1.6) pkg> activate julia/
(julia) pkg> instantiate
(julia) pkg> <BS>
julia> include("julia/copy_data.jl")
BenchmarkTools.Trial: 1 sample with 1 evaluation.
 Single result which took 54.524 s (3.18% GC) to evaluate,
 with a memory estimate of 5.42 GiB, over 107996251 allocations.
```

Development Environment
---------------

Steps to follow if to recreate the tests.

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


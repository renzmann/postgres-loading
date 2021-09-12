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

Python + pandas

```bash
python ./python/copy_using_pandas.py
Finished in 26.1 seconds
```

Julia + LibPQ

```bash
julia
julia> ]
(@v1.6) pkg> activate julia/
(julia) pkg> instantiate
(julia) pkg> <BS>
julia> include("julia/copy_data.jl")

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


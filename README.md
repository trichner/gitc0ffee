# Git Commit Vanity Hash Solver

Neat tool to find a 'vanity' hash for a given git commit. Make all your commits hashes start with the
prefix `c0ffee`, `cafe`, `badc0de5` or whatever makes you happy!

# Install

```bash
go install github.com/trichner/gitc0ffee@latest
```

# Usage

```bash
# do a normal git commit
$ git commit -am '...'

# update the last commit with a vanity hash
$ gitc0ffee --update-ref --prefix c0ffee
```

# Q & A

__Will this break git tooling?__
> Maybe. Not all tooling deals well with prefix collisions. Some tools just deal with short-revisions (7 characters) and
> may therefore break.

__How fast is it?__

- 6 character prefix: less than a second
- 8 character prefix: in the order of one or more minutes

Measured on a __MacBook Pro 16' 2021 with an M1 Max__. Slightly slower
on an __AMD Ryzen 7 5800X__.

__Why not use the GPU?__

Using the GPU is a lot more effort and a lot less portable, since it takes less than a second to brute force
the `c0ffee` prefix there is no need for anything fancier.

The solver implementation can easily be extended though and as
a matter of fact there are already at least three available:

```plain
singlethreaded  - plain Go implementation of single thread brute force
concurrent      - concurrent version of the singlethreaded solver
native          - concurrent solver with hot-loop written in C,
                  slightly faster than 'concurrent' solver

# use:
gitc0ffee --solver <solver> ...
```

__What prefix should I choose?__

All even-length hexadecimal prefix will do (`[0-9a-f]{0,40}`), for cool inspiration
see [Hexspeak](https://en.wikipedia.org/wiki/Hexspeak). Other ideas are repetions or sequences, e.g. `0001`, `0002`, ...

Note that the longer the prefix is, the longer cracking will
take. Prefixes beyond 8 characters may not finish in useful time.

# Implementation Details

Conceptually it roughly works as follows:

1. Get the latest commit digest (`git rev-parse HEAD`).
2. Parse the raw object (`git cat-file -p <digest>`).
3. Add an additional `coffeesalt` header to the commit object and tweak the salt value until a prefix collision is
   found. This is the actual brute-forcing.
4. Write the new commit object to the git store (`git hash-object -t commit --stdin`).
5. (optional) Update the current branch the new commit object (`git update-ref HEAD <new digest>`).

# Previous Work & Inspirations

There are quite a few similar tools. Some are either a bit more on the proof-of-concept side or might need specific GPU
features.

- https://github.com/phadej/git-badc0de
- https://github.com/tochev/git-vanity
- https://github.com/prasmussen/git-vanity-hash
- https://github.com/mattbaker/git-vanity-sha

# Wishlist

- use CPU accelerated assembly, see
  also [Linux Implementation](https://github.com/torvalds/linux/blob/master/arch/x86/crypto/sha1_ni_asm.S)
- GPU accelerated solver - `OpenCL` or `CUDA`


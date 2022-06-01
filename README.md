# Git Commit Vanity Hash Solver

Neat tool to find a 'vanity' hash for a given git commit. E.g. make all your commits hashes start with the
prefix `badc0de` or `deadbeef`.

__Will this break git tooling?__
> Maybe. Not all tooling deals well with prefix collisions. Some tools just deal with short-revisions (7 characters) and
> may therefore break.

# Usage

```bash
# do a normal git commit
$ git commit -am '...'

# update the commit with a vanity hash
$ ./gitc0ffe
c0ffeeb3db6b8b7f72657b3ede42429800b70282

# update the HEAD branch and point to the vanity hash
$ git update-ref HEAD c0ffeeb3db6b8b7f72657b3ede42429800b70282
```

# Previous Work & Inspirations

- https://github.com/phadej/git-badc0de
- https://github.com/tochev/git-vanity
- https://github.com/prasmussen/git-vanity-hash
- https://github.com/mattbaker/git-vanity-sha

# Wishlist

- Use `git hash-object -t commit --stdin` to update the commit object instead of directly writing it.
- add option to also update branch (`update-ref`)


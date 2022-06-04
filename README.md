# Git Commit Vanity Hash Solver

Neat tool to find a 'vanity' hash for a given git commit. E.g. make all your commits hashes start with the
prefix `c0ffee`, `badc0de` or `deadbeef`.

__Will this break git tooling?__
> Maybe. Not all tooling deals well with prefix collisions. Some tools just deal with short-revisions (7 characters) and
> may therefore break.

__How fast is it?__

- 6 character prefix: less than a second
- 8 character prefix: in the order of one or more minutes

Measured on a MacBook Pro 16' 2021 with an M1 Max.

__Why not use the GPU?__

Using the GPU is a lot more effort and a lot less portable, since it takes less than a second to brute force
the `c0ffee` prefix there is no need for anything fancier.

The solver implementation can easily be extended though.

# Usage

```bash
# do a normal git commit
$ git commit -am '...'

# update the last commit with a vanity hash
$ gitc0ffee --update-ref --prefix c0ffee
```

*For inspiration see [Hexspeak](https://en.wikipedia.org/wiki/Hexspeak).*

# How?

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

- add option to also update branch (`update-ref`)


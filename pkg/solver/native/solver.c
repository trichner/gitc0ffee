

#include "solver.h"

#include <inttypes.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include "sha1.h"

/*
ideas for further optimization:
- CPU acceleration (specialized instruction sets) for sha1 computation
*/

static const uint8_t HEX_ALPHABET[] = "0123456789abcdef";
static const int UINT64_BYTES = 8;
static const int SALT_HEX_LEN = 2 * UINT64_BYTES;

static uint8_t to_hex(uint64_t b);

static void update_salt(uint8_t salt_bytes[], uint64_t salt);

static void print_hash(uint32_t hash[static STATE_LEN]);

static bool has_prefix(uint32_t hash[static STATE_LEN], uint8_t prefix[],
                       size_t prefix_len);

static bool should_print_progress(uint64_t salt);

uint64_t solve(uint8_t raw_bytes[], size_t raw_len, uint8_t prefix[],
               size_t prefix_len, size_t salt_offset, uint64_t salt_start,
               uint64_t salt_end) {
  uint8_t *salt_bytes = &raw_bytes[salt_offset];
  uint32_t hash[STATE_LEN];

  clock_t tick = clock();

  for (uint64_t salt = salt_start; salt < salt_end; salt++) {
    update_salt(salt_bytes, salt);
    sha1_sum(raw_bytes, raw_len, hash);

    if (has_prefix(hash, prefix, prefix_len)) {
      //            fprintf(stderr, "found salt: %016x\n", salt);
      //            print_hash(hash);
      return salt;
    }

    if (should_print_progress(salt)) {
      clock_t tock = clock();
      float d = ((float)tock - tick) / CLOCKS_PER_SEC;
      float rate = (float)(salt - salt_start) / (d * 1000.0);
      fprintf(stderr, "brute forcing at %4.f khash/s\n", rate);
    }
  }

  return ERR_SALTS_EXHAUSTED;
}

static bool should_print_progress(uint64_t salt) {
  return (salt & UINT64_C(0xffffff)) == UINT64_C(0xffffff);
}

static bool has_prefix(uint32_t hash[static STATE_LEN], uint8_t prefix[],
                       size_t prefix_len) {
  uint8_t is_matching = 0;

  uint32_t h;
  size_t hash_pos = 0;
  for (int i = 0; i < prefix_len; i++) {
    if (i % 4 == 0) {
      h = hash[hash_pos];
      hash_pos++;
    }

    uint8_t hash_byte = (h >> 24) & 0xff;
    is_matching |= hash_byte ^ prefix[i];
    h <<= 8;
  }

  return is_matching == 0;
}

static void print_hash(uint32_t hash[static STATE_LEN]) {
  for (int i = 0; i < STATE_LEN; i++) {
    fprintf(stderr, "%08x", hash[i]);
  }
  fprintf(stderr, "\n");
}

static void update_salt(uint8_t salt_bytes[], uint64_t salt) {
  for (int i = 0; i < SALT_HEX_LEN; i++) {
    salt_bytes[SALT_HEX_LEN - 1 - i] = to_hex(salt);
    salt >>= 4;
  }
}

static uint8_t to_hex(uint64_t b) {
  b &= 0x0f;
  return HEX_ALPHABET[b];
}

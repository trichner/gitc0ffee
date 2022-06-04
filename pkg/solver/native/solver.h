
#include <stdint.h>
#include <stdio.h>

#include "sha1.h"

#define ERR_SALTS_EXHAUSTED UINT64_MAX

/**
 * Brute forces a prefix hash collision for given bytes.
 * The passed raw_bytes will be updated with incrementing salt values
 * until one leads to the expected prefix collision. If none matches
 * ERR_SALTS_EXHAUSTED is returned.
 *
 * @param raw_bytes    the bytes to calculate the hash on
 * @param raw_len      length of bytes in raw_bytes
 * @param prefix       the bytes of the prefix that should match
 * @param prefix       length of bytes for the prefix
 * @param salt_offset  the offset in bytes within `raw_bytes` where the salt is
 * located, the salt is expected to be a hex-encoded uint64
 * @param salt_start   the first salt to try
 * @param salt_end     the last salt to try
 *
 * @returns the salt leading to the prefix collision or ERR_SALTS_EXHAUSTED if
 * no collision is found.
 */
uint64_t solve(uint8_t raw_bytes[], size_t raw_len, uint8_t prefix[],
               size_t prefix_len, size_t salt_offset, uint64_t salt_start,
               uint64_t salt_end);

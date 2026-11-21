#!/usr/bin/env python3
"""Compare code/data sections in two Windows PE files; no execution needed."""
import argparse
import hashlib
import json
from pathlib import Path
import struct


def pe_sections(path: Path):
    raw = path.read_bytes()
    if raw[:2] != b'MZ':
        raise ValueError(f'{path}: missing MZ header')
    peoff, = struct.unpack_from('<I', raw, 0x3C)
    if raw[peoff:peoff+4] != b'PE\0\0':
        raise ValueError(f'{path}: invalid PE signature')
    count, = struct.unpack_from('<H', raw, peoff + 6)
    optional_size, = struct.unpack_from('<H', raw, peoff + 20)
    start = peoff + 24 + optional_size
    sections = {}
    for i in range(count):
        p = start + 40 * i
        name = raw[p:p+8].rstrip(b'\x00').decode('ascii')
        virtual_size, virtual_address, size, offset = struct.unpack_from('<IIII', raw, p+8)
        if offset+size > len(raw):
            raise ValueError(f'{path}: {name} extends beyond file')
        sections[name] = {
            'size_bytes': size,
            'sha256': hashlib.sha256(raw[offset:offset+size]).hexdigest(),
        }
    return {
        'filename': path.name,
        'size_bytes': len(raw),
        'sha256': hashlib.sha256(raw).hexdigest(),
        'sections': sections,
    }



if __name__ == '__main__':
    args = argparse.ArgumentParser(description=__doc__)
    args.add_argument('exe_a', type=Path)
    args.add_argument('exe_b', type=Path)
    opt = args.parse_args()
    a, b = pe_sections(opt.exe_a), pe_sections(opt.exe_b)
    result = {'a': {k:v for k,v in a.items() if k!='sections'},
              'b': {k:v for k,v in b.items() if k!='sections'},
              'sections': {}}
    for name in sorted(a['sections'].keys() | b['sections'].keys()):
        left, right = a['sections'].get(name), b['sections'].get(name)
        result['sections'][name] = {
            'identical': left is not None and right is not None and left == right,
            'a': left, 'b': right,
        }
    print(json.dumps(result, indent=2))

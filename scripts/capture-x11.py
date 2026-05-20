#!/usr/bin/env python3
import argparse
import ctypes
import os
import sys
from ctypes import POINTER, c_char_p, c_int, c_uint, c_ulong, c_void_p

from PIL import Image


ZPIXMAP = 2
ALL_PLANES = c_ulong(~0 & ((1 << (ctypes.sizeof(c_ulong) * 8)) - 1))


class XImage(ctypes.Structure):
    _fields_ = [
        ("width", c_int),
        ("height", c_int),
        ("xoffset", c_int),
        ("format", c_int),
        ("data", c_void_p),
        ("byte_order", c_int),
        ("bitmap_unit", c_int),
        ("bitmap_bit_order", c_int),
        ("bitmap_pad", c_int),
        ("depth", c_int),
        ("bytes_per_line", c_int),
        ("bits_per_pixel", c_int),
        ("red_mask", c_ulong),
        ("green_mask", c_ulong),
        ("blue_mask", c_ulong),
    ]


def load_x11():
    lib = ctypes.cdll.LoadLibrary("libX11.so.6")
    lib.XOpenDisplay.argtypes = [c_char_p]
    lib.XOpenDisplay.restype = c_void_p
    lib.XDefaultScreen.argtypes = [c_void_p]
    lib.XDefaultScreen.restype = c_int
    lib.XRootWindow.argtypes = [c_void_p, c_int]
    lib.XRootWindow.restype = c_ulong
    lib.XGetGeometry.argtypes = [
        c_void_p,
        c_ulong,
        POINTER(c_ulong),
        POINTER(c_int),
        POINTER(c_int),
        POINTER(c_uint),
        POINTER(c_uint),
        POINTER(c_uint),
        POINTER(c_uint),
    ]
    lib.XGetGeometry.restype = c_int
    lib.XGetImage.argtypes = [
        c_void_p,
        c_ulong,
        c_int,
        c_int,
        c_uint,
        c_uint,
        c_ulong,
        c_int,
    ]
    lib.XGetImage.restype = POINTER(XImage)
    lib.XDestroyImage.argtypes = [c_void_p]
    lib.XDestroyImage.restype = c_int
    lib.XCloseDisplay.argtypes = [c_void_p]
    lib.XCloseDisplay.restype = c_int
    return lib


def mask_shift(mask):
    if mask == 0:
        return 0
    shift = 0
    while mask & 1 == 0:
        shift += 1
        mask >>= 1
    return shift


def mask_bits(mask):
    return int(mask).bit_count()


def channel(pixel, mask, shift, bits):
    if mask == 0 or bits == 0:
        return 0
    value = (pixel & mask) >> shift
    return int(round(value * 255 / ((1 << bits) - 1)))


def capture(output):
    display_name = os.environ.get("DISPLAY", ":99").encode()
    x11 = load_x11()
    display = x11.XOpenDisplay(display_name)
    if not display:
        raise RuntimeError(f"Cannot open X display {display_name.decode()}")

    image_ptr = None
    try:
        screen = x11.XDefaultScreen(display)
        root = x11.XRootWindow(display, screen)

        root_return = c_ulong()
        x = c_int()
        y = c_int()
        width = c_uint()
        height = c_uint()
        border = c_uint()
        depth = c_uint()
        if not x11.XGetGeometry(
            display,
            root,
            ctypes.byref(root_return),
            ctypes.byref(x),
            ctypes.byref(y),
            ctypes.byref(width),
            ctypes.byref(height),
            ctypes.byref(border),
            ctypes.byref(depth),
        ):
            raise RuntimeError("XGetGeometry failed")

        image_ptr = x11.XGetImage(
            display,
            root,
            0,
            0,
            width.value,
            height.value,
            ALL_PLANES,
            ZPIXMAP,
        )
        if not image_ptr:
            raise RuntimeError("XGetImage failed")

        image = image_ptr.contents
        byte_order = "little" if image.byte_order == 0 else "big"
        bytes_per_pixel = max(1, image.bits_per_pixel // 8)
        buffer_size = image.bytes_per_line * image.height
        if not image.data:
            raise RuntimeError("XImage has no pixel data")
        raw = ctypes.string_at(image.data, buffer_size)

        red_shift = mask_shift(image.red_mask)
        green_shift = mask_shift(image.green_mask)
        blue_shift = mask_shift(image.blue_mask)
        red_bits = mask_bits(image.red_mask)
        green_bits = mask_bits(image.green_mask)
        blue_bits = mask_bits(image.blue_mask)

        rgb = bytearray(image.width * image.height * 3)
        out = 0
        for row in range(image.height):
            row_offset = row * image.bytes_per_line
            for col in range(image.width):
                offset = row_offset + col * bytes_per_pixel
                pixel = int.from_bytes(raw[offset : offset + bytes_per_pixel], byte_order)
                rgb[out] = channel(pixel, image.red_mask, red_shift, red_bits)
                rgb[out + 1] = channel(pixel, image.green_mask, green_shift, green_bits)
                rgb[out + 2] = channel(pixel, image.blue_mask, blue_shift, blue_bits)
                out += 3

        os.makedirs(os.path.dirname(os.path.abspath(output)), exist_ok=True)
        Image.frombytes("RGB", (image.width, image.height), bytes(rgb)).save(output)
    finally:
        if image_ptr:
            x11.XDestroyImage(image_ptr)
        x11.XCloseDisplay(display)


def main():
    parser = argparse.ArgumentParser(description="Capture the current X11 root window.")
    parser.add_argument("output", help="PNG output path")
    args = parser.parse_args()
    try:
        capture(args.output)
    except Exception as exc:
        print(f"capture failed: {exc}", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

Photo Box

![logo](camera-logo.png)

## Features

- API
    - `/upload`
    - `/thumbnail`
- Save origin image
- Generate thumbnail image with size your desire
- JSON response format
- Custom `index.html` homepage
- Redis cached upload result based on photo hash

## API

```
/upload
/thumbnail

    Common query parameters:
        width     int    "max thumbnail width"
        height    int    "max thumbnail width"
        quality   int    "thumbnail quality"
```

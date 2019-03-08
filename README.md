# Photo Box

![logo](camera-logo.png)

## Features

- API
    - `/upload`
    - `/thumbnail`
- Save origin image
  - local disk
  - AWS S3
- Generate thumbnail image with size you desire
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

## [Setup S3 storage](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)

* export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_DEFAULT_REGION`
* export `PHOTOBOX_BUCKET`, default `photobox-develop`

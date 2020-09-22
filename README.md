# himitz

## About

A command line tool to encrypt and decrypt using Google Cloud KMS.

## Preapre

### Setup GCP

* Create Service Acccount
* Create JSON Key and download it.

### Setup Google Cloud KMS

* Create KeyRing
* Create Key

### Install himitz

```sh
wget https://github.com/oreno-tools/himitz/releases/download/latest/himitz_darwin_amd64 -O ~/bin/himitz
chmod +x ~/bin/himitz
```

## Usage

### Encrypting string

```sh
echo 'foo' | env GOOGLE_APPLICATION_CREDENTIALS=credential.json himitz -encrypt \
  -project=xxxxxxxxx -ring=xxxxxxx -key=xxxxxxx
```

The ouput looks like this:

```sh
Encrypted data with base64 encoded: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Decrypting data

```sh
echo 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' | env GOOGLE_APPLICATION_CREDENTIALS=credential.json himitz -decrypt \
  -project=xxxxxxxxx -ring=xxxxxxx -key=xxxxxxx
```

The ouput looks like this:

```sh
Decrypted data: foo
```

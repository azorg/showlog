# Prepare build enviroment

## Install Go
As you wish.

## Install make
```
$ apt install make
```

## Install system dependencies
```
$ sudo apt-get install -y gcc libgl1-mesa-dev xorg-dev
```

## Install Fyne (look fyne.io)
```
$ go install fyne.io/fyne/v2/cmd/fyne@latest
```
Add ~/go/bin to PATH.

## Build for PC
```
$ make
```

## Build for Android (apk)

### Install android-ndk
Set enviroment variable (e.g):
```
$ export ANDROID_NDK_HOME="${HOME}/android-ndk"
```
### Run make
```
$ make apk
```


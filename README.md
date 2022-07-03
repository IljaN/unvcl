# unvcl
A cli tool and lib to extract sound-files from *.vcl files. VCL is a file-format which is mainly used by some early DOS 
games developed by Epic Megagames.  

## Building from source
There are currently no binary releases available so you need to build unvcl from source. Clone this repo then run 
``make build_all`` to build binaries for Linux, MacOS and Windows.

## Usage 
unvcl expects two arguments: The source vcl file and a target directory in to which the samples should be extracted.

### Example
Extract all PCM samples from jill1.vcl and store them as WAV file in the current directory
```shell 
$ unvcl ~/Games/JILL1/jill1.vcl . 
```
List extracted files
```shell
$ ls *.wav
jill1_0.wav   jill1_11.wav  jill1_13.wav  jill1_15.wav  jill1_17.wav  jill1_19.wav  jill1_20.wav  jill1_22.wav  jill1_3.wav  jill1_5.wav  jill1_7.wav  jill1_9.wav
jill1_10.wav  jill1_12.wav  jill1_14.wav  jill1_16.wav  jill1_18.wav  jill1_1.wav   jill1_21.wav  jill1_2.wav   jill1_4.wav  jill1_6.wav  jill1_8.wav
```
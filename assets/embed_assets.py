#This embeds all of the assets into a .go file 
#by zipping the contents and then encoding them as a string constant

import gzip
import os
import io

#Read in all the file data and compress
data = {}
for root, dirs, files in os.walk(".", topdown=True):
  for file_name in files:
    if file_name.endswith((".png", ".wav", ".ogg")):
      with open(file_name, "rb") as fin:
        root_name, extension = os.path.splitext(file_name)
        key_name = extension.upper().strip(".") + "_" + root_name.upper()
        #stream = io.BytesIO(fin.read())
        data[key_name] = gzip.compress(fin.read())
        print("Read", file_name, "into", key_name)

#Write as a string constant in a .go file
with open("assets.go", "w", encoding='utf-8') as fout:
  fout.write("package assets\n")
  for key in data:
    fout.write("const " + key + '="')
    for byte in data[key]:
      #Each byte is offset so that it may be printed as a single valid unicode character
      byte += 186
      fout.write(chr(byte))
    fout.write('"\n')
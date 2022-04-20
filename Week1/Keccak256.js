from base64 import encode
import sha3
import os

print("Keccak 256 di python\n")
namaorang = input("Nama Orang: ")
os.system('CLS')
print("Nama Orang: \n", namaorang)
encoded = namaorang.encode()
obj_encoded = sha3.keccak_256(encoded)
print("Nama orang sesudah hash Keccak 256: \n", obj_encoded.hexdigest())
